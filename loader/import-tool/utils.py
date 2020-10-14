import csv
import numpy as np
import math
from scipy.cluster.hierarchy import dendrogram, linkage
from scipy.spatial.distance import pdist
import scipy.spatial.distance as ssd
import matplotlib.pyplot as plt
import pandas as pd

"""
    Support method for populating the patient to concepts matrix.
    Putting a 1 in the cell identified by a certain (patient_num, concept_cd) pair.
"""
def fill_patient_to_concept_table(patient_to_concept, observation_fact_df):
    observation_fact_df['patient_num'] = observation_fact_df['patient_num'].astype(str)
    for row in observation_fact_df.itertuples():
        patient_to_concept[row[3]][row[2]] = 1

def permutation_entropy(clusters, needed_dummies_per_cluster):
    entropy = 0
    for cluster in clusters:
        entropy = entropy + math.log(np.math.factorial(cluster.shape[0]))
        
    return entropy

"""
    Method for plotting the histograms passed as parameters
"""
def plot_histograms(original_histogram=None, theoretical_dummies=None, histogram_with_dummies=None):
    plt.figure(figsize=(20,15))
    plt.subplot(212)
    plt.xlabel('concept_cd', fontsize=15)
    plt.ylabel('number of patients', fontsize=15)

    to_handle = []

    if (original_histogram is not None):
        original, = plt.plot(np.arange(original_histogram.shape[0]), original_histogram, c='blue', label='Original Histogram')
        to_handle.append(original)
    
    if (theoretical_dummies is not None):
        theoretical, = plt.plot(np.arange(theoretical_dummies.shape[0]), theoretical_dummies, c='red', label='Theoretical Dummies')
        to_handle.append(theoretical)

    if (histogram_with_dummies is not None):
        real, = plt.plot(np.arange(histogram_with_dummies.shape[0]), histogram_with_dummies, c='green', label='Real Dummies')
        to_handle.append(real)

    plt.legend(handles=to_handle, fontsize='large')

    plt.show()
    
"""
    Support method for getting the size of the smallest cluster
"""
def get_smallest_anonymity_set(cluster_labels):
    smallest_anonymity_set = cluster_labels[cluster_labels == 0].size
    
    for cluster_label in cluster_labels:
        anonymity_set = cluster_labels[cluster_labels == cluster_label].size
        if smallest_anonymity_set > anonymity_set:
            smallest_anonymity_set = anonymity_set
            
    return smallest_anonymity_set    
    
"""
    Support method for getting the biggest cluster
"""
def get_biggest_cluster(cluster_labels):
    biggest_cluster = (cluster_labels[cluster_labels == 0], cluster_labels[cluster_labels == 0].size)
    
    for cluster_label in cluster_labels:
        cluster = cluster_labels[cluster_labels == cluster_label]
        anonymity_set = cluster.size
        if anonymity_set > biggest_cluster[1]:
            biggest_cluster = (cluster, cluster.size)
    
    return biggest_cluster[0]

"""
    Returns a partitioned histogram and the corresponding one with the dummies added
"""
def build_clustered_histograms(histogram, cluster_labels):
    new_histogram = pd.Series([])
    new_full_histogram = pd.Series([])
    
    for cluster_label in np.unique(cluster_labels):
        cluster = cluster_labels[cluster_labels == cluster_label]
        cluster = histogram[cluster.index]
        new_histogram = pd.concat([new_histogram, cluster])
        
        needed_dummies = np.max(cluster) - cluster
        
        new_full_histogram = pd.concat([new_full_histogram, cluster + needed_dummies])
        
    return new_histogram, new_full_histogram

"""
    Assigning the labels computed by kmeans to all the codes belonging to a certain hierarchical cluster.
    Different hierarchical clusters are merged
"""
def assign_labels_to_original_codes(new_clusters, clusters_hierarchical):
    original_codes_clusters = clusters_hierarchical.copy()
    for label in np.unique(new_clusters):
        symbolic_codes = new_clusters[new_clusters == label]
        #clusters_hierarchical_labels = symbolic_codes.index.map(lambda x: __get_symbolic_code_original_label(x))
        clusters_hierarchical_labels = symbolic_codes.index.map(lambda x: x.split(':')[1])
        
        for original_label in np.unique(clusters_hierarchical_labels):
            original_cluster = clusters_hierarchical[clusters_hierarchical == int(original_label)]
            original_codes_clusters[original_cluster.index] = label
                              
    return original_codes_clusters
