import random
import pandas as pd
import utils

# returns the name of a patient that will be duplicated in order to add occurrences to codes in a cluster
"""
    Getting the name of the patient that will be duplicated in order to add occurrences to codes in a cluster
"""
def get_patient(patient_concepts_matrix, number_codes_needing_patients, codes_in_cluster, last, submatrix):
    # take only the columns that relate to the codes in the cluster of interest
    patient_number_codes = submatrix.sum(1)
    patients_with_same_number_codes = patient_number_codes[(patient_number_codes == number_codes_needing_patients) & (patient_number_codes > 0)]
    if(patients_with_same_number_codes.shape[0] > 0):
        #print("same number of codes")
        patient_index = random.randint(0, patients_with_same_number_codes.shape[0] - 1)
        return patients_with_same_number_codes.index[patient_index]
    
    patients_with_less_number_codes = patient_number_codes[(patient_number_codes < number_codes_needing_patients) & (patient_number_codes > 0)]
    if (patients_with_less_number_codes.shape[0] > 0):
        #print("smaller number of codes")
        # get the patients that have the maximum number of codes
        candidate_patients = patients_with_less_number_codes[patients_with_less_number_codes == patients_with_less_number_codes.max(0)]
        patient_index = random.randint(0, candidate_patients.shape[0] - 1)
        patient = candidate_patients.index[patient_index]
        
        return patient
    
    patients_with_more_number_codes = patient_number_codes[(patient_number_codes > number_codes_needing_patients) & (patient_number_codes > 0)]
    if(patients_with_more_number_codes.shape[0] > 0):
        #print("patients with more codes than needed: ")
        candidate_patients = patients_with_more_number_codes[patients_with_more_number_codes == patients_with_more_number_codes.min(0)]
        patient_index = random.randint(0, candidate_patients.shape[0] - 1)
        
        patient = candidate_patients.index[patient_index]
        
        return patient
    
    raise ValueError("Problem in retrieving a patient to duplicate")
    return 0

"""
    Checking that the generated dummy is not a duplicate in any cluster
"""
def check_duplicate_all_clusters(patient_concepts_matrix, dummy_patient, clusters_label):
    # taking only the rows with the same weight of the dummy patient
    filtered_patient_concepts_matrix = patient_concepts_matrix[patient_concepts_matrix.sum(1) == dummy_patient.sum()]

    for label in clusters_label.unique():
        if(check_duplicate_single_cluster(filtered_patient_concepts_matrix[clusters_label[clusters_label == label].index], dummy_patient[clusters_label[clusters_label == label].index])):
            return True
        
    return False
    
"""
    Checking that the generated dummy is not a duplicate in a specific cluster
"""
def check_duplicate_single_cluster(patient_concepts_matrix_cluster, x_cluster):
    if ((x_cluster - patient_concepts_matrix_cluster == 0).all(1).any()):
        return True
        
    return False

def fix_dummy_duplicate(histogram, dummy_patient, clusters_labels, patient_concepts_matrix):
    #print("FIXING DUMMY PATIENT! ***************** : ", dummy_patient.name)
    labels = clusters_labels.unique().copy()
    random.shuffle(labels)

    # considering only the patients with the same weight as the dummy
    filtered_patient_concepts_matrix = patient_concepts_matrix[patient_concepts_matrix.sum(1) == dummy_patient.sum()]
    
    for label in labels:
        # contains #occurrences of codes
        current_cluster = histogram[clusters_labels[clusters_labels == label].index]
        # contains 1 or 0 in correspondence of codes
        patient_cluster = dummy_patient[current_cluster.index].copy()
        
        # takes histogram values and sort them
        codes_set_to_1 = current_cluster[patient_cluster == 1].sort_values(ascending=False).index
        codes_set_to_0 = current_cluster[patient_cluster == 0].sort_values(ascending=True).index
        
        if (codes_set_to_1.shape[0] == 0 or codes_set_to_0.shape[0] == 0):
            #raise ValueError("The patient that was replicated for dummy generation has a cluster with all the codes or a cluster with no code. Perform the clustering again with a greater minimum anonymity set")
            continue

        iterations = 0
        # try to slightly change the dummy patient
        for c_1 in codes_set_to_1:
            for c_0 in codes_set_to_0:
                x = patient_cluster.copy()
                x[c_1] = 0
                x[c_0] = 1

                if (iterations > 200):
                    #raise ValueError("Too many iterations for fixing the dummy ...")
                    break

                # the matrix could (and probably should) be reduced before the check!
                duplicate = check_duplicate_single_cluster(filtered_patient_concepts_matrix[clusters_labels[clusters_labels == label].index], x[clusters_labels[clusters_labels == label].index])
                
#                 duplicate = False
                if not duplicate:
                    # good, we can keep this dummy
                    break
                iterations = iterations + 1
                
            if (not duplicate or iterations > 200):
                # good, let's keep this dummy
                break
        #if (duplicate):
            #raise ValueError("This is a duplicate patient. Try to increase the anonymity set")
        
        dummy_patient[x.index] = x

    #print("FINISHED FIXING THE DUMMY! *****************************")
    
    return dummy_patient

"""
    This method is used to create the dummy patient, by shuffling the codes
    of the patient that was chosen for duplication.
"""
def perform_shuffle(patient_row, clustered_histogram, histogram_dummies, clusters_labels):
    dummy_patient = patient_row.copy()
    for label in clusters_labels.unique():
        # get row portion that relates to the cluster
        current_cluster = clustered_histogram[clusters_labels[clusters_labels == label].index]
        current_cluster_dummies = histogram_dummies[current_cluster.index]
        patient_cluster = dummy_patient[current_cluster.index]
        
        codes_needing_patients_cluster = current_cluster[(current_cluster_dummies - current_cluster).abs() >= 1]
        number_codes_needing_patients_cluster = codes_needing_patients_cluster.shape[0]

        k = patient_cluster.sum(0)
        if (k <= number_codes_needing_patients_cluster):
            # take k random codes that need more occurrences
            k_random_elements = pd.Index(random.sample(list(codes_needing_patients_cluster.index), k))
            patient_cluster[patient_cluster.index] = 0
            patient_cluster[k_random_elements] = 1
        else:
            # always include all the codes that need occurrences
            # also include the codes that already have reached the need number of occurrences (but take the ones that appear less often)
            k_indices = current_cluster.values.argpartition(k - 1)
            k_min_codes = current_cluster[k_indices[:k]].index

            patient_cluster[patient_cluster.index] = 0
            patient_cluster[k_min_codes] = 1

        dummy_patient[patient_cluster.index] = patient_cluster[patient_cluster.index]

    return dummy_patient

def add_dummy_patient(clusters_labels, clustered_histogram, histogram_with_dummies, patient_concepts_matrix, dummy_n, label, mapping, last, submatrix):
    # need to choose the patient to duplicate

    # taking the maximum number of occurrences of a cluster
    # the elements of the cluster, along with the label
    cluster_elements = clusters_labels[clusters_labels == label]
    
    # taking all the codes within that cluster, along with the maximum occurrences (they all have the same value)
    filled_cluster = histogram_with_dummies[cluster_elements.index]

    histogram_cluster = clustered_histogram[cluster_elements.index]
    codes_needing_patients = histogram_cluster[histogram_cluster < filled_cluster[0]]
    
    # number of codes that need more occurrences in THIS cluster
    number_codes_needing_patients = codes_needing_patients.shape[0]
    # need to get a patient with this number of codes in this cluster

    if (number_codes_needing_patients == 0):
        return patient_concepts_matrix, mapping

    #patient = get_patient(patient_concepts_matrix, number_codes_needing_patients, filled_cluster, last)
    patient = get_patient(patient_concepts_matrix, number_codes_needing_patients, filled_cluster, last, submatrix)

    new_row = perform_shuffle(patient_concepts_matrix.loc[patient], clustered_histogram, histogram_with_dummies, clusters_labels)
    new_row.name = str(int(last) + dummy_n)

    new_patient_concepts_matrix = add_dummy_to_patient_concepts(patient_concepts_matrix, new_row, clustered_histogram, clusters_labels)
    mapping[new_row.name] = patient
    
    return new_patient_concepts_matrix, mapping

"""
    This method adds the dummy patient (new_row) to the patient_concepts matrix.
    It does not check whether the new generated row is duplicated.
"""
def add_dummy_to_patient_concepts(patient_concepts_matrix, new_row, histogram, clusters_labels):
    # perform checks on this new_row, in order to see if it is equal to any existing patient
    # in case: fix it

    fixed_dummy = new_row

    #if(check_duplicate_all_clusters(patient_concepts_matrix, new_row, clusters_labels)):
    #    fixed_dummy = fix_dummy_duplicate(histogram, new_row, clusters_labels, patient_concepts_matrix)
    #else:
    #    fixed_dummy = new_row
        
    df = pd.DataFrame(fixed_dummy).T
    patient_concepts_matrix = patient_concepts_matrix.append(df)
    
    return patient_concepts_matrix

"""
    This function returns a patient to concepts matrix that contains the dummy that relates to each patient.
    The returned matrix has twice the number of rows of the input matrix.
"""
def duplicate_patients(patient_concepts_matrix, clusters_labels, dummy_n, mapping, last):
    histogram = patient_concepts_matrix.sum(0)
    clustered_histogram, histogram_with_dummies = utils.build_clustered_histograms(histogram, clusters_labels)

    for patient in patient_concepts_matrix.index:
        histogram = patient_concepts_matrix.sum(0)
        new_row = perform_shuffle(patient_concepts_matrix.loc[patient], histogram, histogram_with_dummies, clusters_labels)
        new_row.name = str(int(last) + dummy_n)

        patient_concepts_matrix = add_dummy_to_patient_concepts(patient_concepts_matrix, new_row, histogram, clusters_labels)
        mapping[new_row.name] = patient
    
        dummy_n = dummy_n + 1
        
    return patient_concepts_matrix, dummy_n, mapping

"""
    This function fills a cluster and returns the updated patient to concepts matrix
"""
def fill_bucket(clusters_labels, patient_concepts_matrix, label, dummy_n, mapping, last):
    labels_right_cluster = clusters_labels[clusters_labels == label]
    clustered_histogram, histogram_filled = utils.build_clustered_histograms(patient_concepts_matrix.sum(0), clusters_labels)

    condition = ((clustered_histogram[labels_right_cluster.index] - histogram_filled[labels_right_cluster.index]).abs() > 1).any() 

    submatrix = patient_concepts_matrix[labels_right_cluster.index]
    submatrix = submatrix[submatrix.index.astype(int) <= int(last)]
    iterations = 0
    max_iterations = 1000 # this number depends on the dataset
    while(condition and iterations < max_iterations):
        patient_concepts_matrix, mapping = add_dummy_patient(clusters_labels, clustered_histogram, histogram_filled, patient_concepts_matrix, dummy_n, label, mapping, last, submatrix)
        clustered_histogram, histogram_filled = utils.build_clustered_histograms(patient_concepts_matrix.sum(0), clusters_labels)
        
        condition = ((clustered_histogram[labels_right_cluster.index] - histogram_filled[labels_right_cluster.index]).abs() > 1).any()
        dummy_n = dummy_n + 1
        iterations = iterations + 1
      
    print("Label ", label, " terminated. Iterations: ", iterations)
    #print("**************** condition: ", condition)
    return patient_concepts_matrix[clustered_histogram.index], dummy_n, mapping
