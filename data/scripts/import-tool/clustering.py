import pandas as pd
from utils import fill_patient_to_concept_table
import scipy.spatial.distance as ssd
import utils
from sklearn.cluster import KMeans
import warnings
import numpy as np
from dummygen import duplicate_patients, fill_bucket
from scipy.spatial.distance import pdist
from scipy.cluster.hierarchy import linkage, fcluster, dendrogram
import generate_csv

""" Implementation of the class that keeps the state of the various steps of the algorithms """

class HierarchicalClustering():
    def __init__(self, observations_path, concept_modifier_separator=""):
        # initialization of structures

        # structures that are filled when the HAC linkage is performed
        self.condensed_dissimilarity_matrix = None
        self.extended_dissimilarity_matrix = None
        self.linkage_matrix = None
        # the linkage needs to be performed before the clustering
        self.linkage = False

        # contains the labels after performing the second phase of the clustering
        self.labels = None
        # the clustering needs to be performed before the generation of dummies
        self.clustering = False

        # structures that are filled when the dummies are generated
        self.patient_concepts_matrix_dummies = None
        self.dummy_to_patient = None
        # the dummies need to be generated before the generation of the new csv files
        self.dummies = False


        # reading the observation fact csv file
        observation_fact_df = pd.read_csv(observations_path, low_memory=False)
        #observation_fact_df.loc[observation_fact_df['modifier_cd'] != '@', 'concept_cd'] = observation_fact_df.loc[observation_fact_df['modifier_cd'] != '@', 'concept_cd'] + concept_modifier_separator + observation_fact_df.loc[observation_fact_df['modifier_cd'] != '@', 'modifier_cd']
        self.observation_fact_df = observation_fact_df

        # renaming any field to 'counter'
        grouped_data = observation_fact_df.groupby(['patient_num', 'concept_cd']).count().rename(columns={'encounter_num':'counter'})

        patient_nums = grouped_data.index.levels[0].astype(str)
        concept_codes = grouped_data.index.levels[1]
        histogram = grouped_data.groupby('concept_cd').count()['counter']
        # discard codes that appear in all patients
        histogram = histogram[histogram < patient_nums.shape[0]].sort_values()

        patient_concepts_matrix = pd.DataFrame(columns=list(concept_codes), index=list(patient_nums)).fillna(0)


        fill_patient_to_concept_table(patient_concepts_matrix, observation_fact_df)
        

        # discard codes that appear in all patients
        patient_concepts_matrix = patient_concepts_matrix[histogram.index]

        self.patient_concepts_matrix = patient_concepts_matrix
        self.histogram = histogram

    """
        This function builds the linkage matrix: it defines clusters and the distance between them (the hierarchical clustering is applied)
        distance_metric: see https://docs.scipy.org/doc/scipy/reference/generated/scipy.spatial.distance.pdist.html
        linkage_method: see https://docs.scipy.org/doc/scipy/reference/generated/scipy.cluster.hierarchy.linkage.html#scipy.cluster.hierarchy.linkage
    """
    def perform_HAC_linkage(self, distance_metric='jaccard', linkage_method='average'):
        if(distance_metric == 'correlation'):
            condensed_dissimilarity_matrix = ssd.squareform(1 - self.patient_concepts_matrix.corr().abs())
        elif(distance_metric == 'personalized_similarity'):
            condensed_dissimilarity_matrix = ssd.squareform(1 - utils.personalized_similarity(self.patient_concepts_matrix))
        else:
            condensed_dissimilarity_matrix = pdist(self.patient_concepts_matrix.T, distance_metric)
    
        Z = linkage(condensed_dissimilarity_matrix, linkage_method)
                
        self.set_dissimilarity_matrix(condensed_dissimilarity_matrix)
        self.linkage_matrix = Z
        
        self.linkage = True

        return self.linkage_matrix, condensed_dissimilarity_matrix

    def set_dissimilarity_matrix(self, condensed_dissimilarity_matrix):
        # dissimilarity matrix stored as a vector
        self.condensed_dissimilarity_matrix = condensed_dissimilarity_matrix

        # dissimilarity matrix stored as a matrix
        self.extended_dissimilarity_matrix = pd.DataFrame(index=self.patient_concepts_matrix.columns, columns=self.patient_concepts_matrix.columns, data=ssd.squareform(condensed_dissimilarity_matrix))
        
    """
        This function define the clusters of codes that will be considered for the generation of dummies.
        It consists in defining further clusters, besides the ones computed through the HAC.
        It has to be called after the linkage. Since the linkage matrix is not touched through the process, it is possible to call
        this method multiple times (results may vary among executions: only the last one will be kept).
        min_anonymity_set_allowed: defines the minimum allowed size of the smallest cluster
        max_ratio: defines how much bigger the histogram with dummies is, in comparison with the histogram without dummies
        HAC_parameters: define the threshold for cutting the clusters from the HAC and the criterion that is used for doing it. Also see https://docs.scipy.org/doc/scipy/reference/generated/scipy.cluster.hierarchy.fcluster.html#scipy.cluster.hierarchy.fcluster
    """
    def perform_clustering(self, min_anonymity_set_allowed=1, max_ratio=10, HAC_parameters=None, n_clusters=3):
        if (not self.linkage):
            raise ValueError('Linkage not performed')

        # default parameters
        if (HAC_parameters == None):
            HAC_parameters = {}
            HAC_parameters['threshold'] = 0.4
            HAC_parameters['clustering_criterion'] = 'distance'
            
        HAC_labels = fcluster(self.linkage_matrix, HAC_parameters['threshold'], criterion=HAC_parameters['clustering_criterion'])
        HAC_clusters = pd.Series(index=self.patient_concepts_matrix.columns, data=HAC_labels)

        self.labels = self.__perform_clustering_aux(min_anonymity_set_allowed, max_ratio, n_clusters, HAC_clusters)
        self.clustering = True

    # Support method for defining clusters that guarantee a minimum anonymity set
    def __perform_clustering_aux(self, min_anonymity_set_allowed, max_ratio, n_clusters, HAC_clusters):
        labels = self.__apply_k_means(HAC_clusters, self.histogram, k=n_clusters)
        
        smallest_anonymity_set = utils.get_smallest_anonymity_set(labels)
        _, filled_histogram = utils.build_clustered_histograms(self.histogram, labels)
        ratio = filled_histogram.sum() / self.histogram.sum()

        
        if (smallest_anonymity_set < min_anonymity_set_allowed and ratio > max_ratio):
            warnings.warn("Anonymity set below the minimum one and ratio greater than the maximum! Try with a greater minimum anonymity set and/or a greater maximum ratio")
            return labels
        
        if (ratio > max_ratio):
            # having more clusters implies a more precise clustering and a smaller ratio. However, the minimum anonymity set will decrease
            return self.__perform_clustering_aux(min_anonymity_set_allowed, max_ratio, n_clusters+1)
            
        if (smallest_anonymity_set < min_anonymity_set_allowed):
            # to adapt the smallest anonimity set to the minimum allowed one. This may cause the ratio to increase noticeably.
            return self.__fill_smallest_buckets(labels, min_anonymity_set_allowed)
        
        return labels

    """
        For performing K means clustering after HAC is applied. This step guarantees that a fixed number of clusters is obtained.
    """
    def __apply_k_means(self, clusters_hierarchical, histogram, k=3):
        kmeans = KMeans(n_clusters=k)
        
        #symbolic_codes = utils.get_symbolic_codes_from_HAC(histogram, clusters_hierarchical)
        # Preparing data for kmeans clustering: passing points in a single hierarchical cluster as a single point (with the mean value)
        symbolic_codes = pd.Series()
        for label in np.unique(clusters_hierarchical):
            cluster = clusters_hierarchical[clusters_hierarchical == label]
            #symbolic_code = __get_symbolic_code_name_from_label(label)
            symbolic_code = 'CLUSTER:' + str(label)
            symbolic_codes[symbolic_code] = round(histogram[cluster.index].mean())

        #return symbolic_codes

        symbolic_codes_reshaped = pd.DataFrame(columns=['HAC_label'], index=symbolic_codes.index, data=symbolic_codes.values.reshape(-1, 1))
    
        kmeans.fit(symbolic_codes_reshaped)
        new_clusters = pd.Series(index=symbolic_codes.index, data=kmeans.predict(symbolic_codes_reshaped))
        
        # Assigning the labels computed by the kmeans to all the codes belonging to a certain hierarchical cluster
        final_labels = utils.assign_labels_to_original_codes(new_clusters, clusters_hierarchical)

        
        return final_labels

    """
        For filling the smallest cluster, in order to meet the requirements concerning the minimum allowed anonymity set.
    """
    def __fill_smallest_buckets(self, clusters, min_anonymity_set_allowed):
        clusters_copy = clusters.copy()
        for cluster_label in np.unique(clusters_copy):
            cluster = clusters_copy[clusters_copy == cluster_label]
            to_add =  min_anonymity_set_allowed - cluster.size
            if to_add > 0:
                # assuming that elements inside the biggest cluster are not correlated.. this assumption may not always be valid
                # and this piece of code may be adapted to behave differently if some different dataset is used
                biggest_cluster = utils.get_biggest_cluster(clusters_copy)
                if to_add >= biggest_cluster.size or biggest_cluster.size - to_add < min_anonymity_set_allowed:
                    warnings.warn("Not enough elements for getting the required anonymity set (try to relax the constraints concerning the required minimum anonymity set and maximum ratio)")
                    return clusters_copy
                clusters_copy[biggest_cluster[np.arange(0, to_add)].index] = cluster_label
                
        return clusters_copy
    
    """
        For plotting the histogram of the codes' occurrences, sorted by number of occurrences
    """
    def plot_histogram(self):
        utils.plot_histograms(original_histogram=self.histogram)
        
    """
        For plotting the histogram of the codes' occurrences, along with the histogram filled with dummy patients' occurrences on codes.
        This is a theoretical plot: it shows the ideal number of occurrences, without really considering the actual process of dummy generation.
        Plot divided in clusters.
    """
    def plot_histogram_with_dummies_theoretical(self):
        if not self.clustering:
            raise ValueError('Clustering not performed yet')

        new_histogram, full_histogram = utils.build_clustered_histograms(self.histogram, self.labels)
        utils.plot_histograms(original_histogram=new_histogram, theoretical_dummies=full_histogram)

    """
        For plotting the histogram of codes' occurrences, along with the histogram filled with dummy patients' occurrences on codes (both theoretical and real ones).
        This plot also considers the actual process of dummy generation'
    """
    def plot_histogram_with_dummies_real(self):
        if not self.dummies:
            raise ValueError('Dummies not generated yet')

        _, theoretical_dummies = utils.build_clustered_histograms(self.histogram, self.labels)
        histogram_with_dummies = self.patient_concepts_matrix_dummies.sum(0)

        cluster_indexing = theoretical_dummies.index

        utils.plot_histograms(self.histogram[cluster_indexing], theoretical_dummies[cluster_indexing], histogram_with_dummies[cluster_indexing])

    # a method for increasing the anonymity set. Not used, but can be considered as an alternative approach.
    def increase_anonymity_set(self, parameters):
        clustering_criterion = parameters['clustering_criterion']
        threshold = parameters['threshold']
        
        min_distance = self.linkage_matrix[:, 2].min()
        max_distance = self.linkage_matrix[:, 2].max()
        offset = (max_distance - min_distance) / 500
        
        if clustering_criterion == 'distance' or clustering_criterion == 'inconsistent':
            return threshold + offset
        else:
            raise NotImplementedError("Other clustering criteria need to be considered")

    # a method for decreasing the ratio. Not used, but can be considered as an alternative approach.
    def decrease_ratio(self, parameters):
        clustering_criterion = parameters['clustering_criterion']
        threshold = parameters['threshold']
        
        min_distance = self.linkage_matrix[:, 2].min()
        max_distance = self.linkage_matrix[:, 2].max()
        offset = (max_distance - min_distance) / 500
        
        if clustering_criterion == 'distance':
            return threshold - offset
        else:
            raise NotImplementedError("Other clustering criteria need to be considered")

    """
        This method visualizes the hierarchical clustering
    """
    def plot_dendrogram(self, granularity=None):
        plt.figure(figsize=(25, 10))
        plt.title('Hierarchical Clustering Dendrogram (truncated)')
        plt.xlabel('concept_cd')
        plt.ylabel('distance')

        if granularity is None:
            granularity = self.patient_concepts_matrix.shape[1]

        dendrogram(
            self.linkage_matrix,
            truncate_mode='lastp',  # show only the last p merged clusters
            p=granularity,  # show only the last p merged clusters
            show_leaf_counts=False,  # otherwise numbers in brackets are counts
            leaf_rotation=90.,
            leaf_font_size=12.,
            show_contracted=True,  # to get a distribution impression in truncated branches
            labels=self.patient_concepts_matrix.columns
        )
        plt.show()

    """
        This method generates the dummies and populates the related data structures to contain this information.
        ascending: flag indicating whether the clusters are filled starting from the one with the smallest higher frequency
    """
    def generate_dummies(self, ascending=True):
        if not self.clustering:
            raise ValueError('Clustering not performed yet')

        # initialize values for dummy generation
        dummy_n = 1
        clusters_labels = self.labels
        _, filled_histogram = utils.build_clustered_histograms(self.histogram, self.labels)
        labels_iterate = self.labels[filled_histogram.sort_values(ascending=ascending).index].unique()

        new_patient_concepts_matrix = self.patient_concepts_matrix.copy()
        mapping = pd.Series({})
        last = new_patient_concepts_matrix.index.astype(int).max() + 1

        new_patient_concepts_matrix, dummy_n, mapping = duplicate_patients(new_patient_concepts_matrix, clusters_labels, dummy_n, mapping, last)

        print("Water filling clusters' labels order: ", labels_iterate)
        for label in labels_iterate:
            new_patient_concepts_matrix, dummy_n, mapping = fill_bucket(clusters_labels, new_patient_concepts_matrix, label, dummy_n, mapping, last)

        self.patient_concepts_matrix_dummies = new_patient_concepts_matrix
        self.dummy_to_patient = mapping
        self.dummies = True

        return self.patient_concepts_matrix_dummies, self.dummy_to_patient

    """
        For getting the theoretical ratio: it is computed by considering the histogram filled up to the maximum value of the histogram
    """
    def theoretical_ratio(self):
        if not self.clustering:
            raise ValueError('Clustering still needs to be performed')

        _, filled_histogram = utils.build_clustered_histograms(self.histogram, self.labels)

        return filled_histogram.sum() / self.histogram.sum()

    """
        For getting the real ratio: it is computed by considering the histogram filled after
        adding all the necessary patients according to the filling criterion for the clusters.
        It will probably be bigger than the theoretical ratio, since more patients are added in order to flatten the clusters.
    """
    def real_ratio(self):
        if not self.dummies:
            raise ValueError('Dummies not generated yet')

        return self.patient_concepts_matrix_dummies.sum(0).sum() / self.histogram.sum()

    """
        For generating the new patient_dimension.csv that includes the dummies
    """
    def generate_patient_dimension(self, original_patient_dimension_path, new_patient_dimension_path):
        if not self.dummies:
            raise ValueError('Dummies not generated yet')

        generate_csv.generate_patient_dimension(self.dummy_to_patient, original_patient_dimension_path, new_patient_dimension_path)

    """
        For generating the new observation_fact.csv that includes the dummies
    """
    def generate_observation_fact(self, new_observation_fact_path):
        if not self.dummies:
            raise ValueError('Dummies not generated yet')

        generate_csv.generate_observation_fact(self.observation_fact_df, self.histogram, self.patient_concepts_matrix_dummies, self.patient_concepts_matrix, self.labels, new_observation_fact_path)
