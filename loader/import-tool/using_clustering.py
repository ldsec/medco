#!/usr/bin/env python
# coding: utf-8

# In[1]:

from clustering import HierarchicalClustering
import csv
import pandas as pd
import sys

original_observation_fact = "data/original/observation_fact.csv"
original_patient_dimension = "data/original/patient_dimension.csv"

output_observation_fact = "data/output/observation_fact.csv"
output_patient_dimension = "data/output/patient_dimension.csv"
output_dummy_to_patient = "data/output/dummy_to_patient.csv"

if len(sys.argv) == 6:
    original_observation_fact = sys.argv[1]
    original_patient_dimension = sys.argv[2]
    output_observation_fact = sys.argv[3]
    output_patient_dimension = sys.argv[4]
    output_dummy_to_patient = sys.argv[5]

print("Started hierarchical clustering.")

# Creating an instance of the clustering wrapper.
# Need to specify the observation fact path and (optionally) the separator between concept_cd 
# and modifier_cd in the new observation_fact.csv file
clustering_wrapper = HierarchicalClustering(original_observation_fact, concept_modifier_separator="_")

# performing hiearchical clustering
# need to specify the similarity metric and linkage method
clustering_wrapper.perform_HAC_linkage('jaccard', 'average')

# performing k-means and adapting the minimum anonymity set requirements
# need to specify the maximum ratio, the minimum allowed anonymity set and the number of clusters
clustering_wrapper.perform_clustering(max_ratio=20, min_anonymity_set_allowed=1000, n_clusters=2)

print("Finished hierarchical clustering.")

print("Started dummy generation.")

# performing dummy generation
patient_concepts_matrix_dummies, dummy_to_patient = clustering_wrapper.generate_dummies(True)


# generating the new patient dimension
# need to specify the old patent_dimension.csv path and the new one
clustering_wrapper.generate_patient_dimension(original_patient_dimension, output_patient_dimension)


# generating the new observation fact
# need to specify the old observation_fact.csv and the new one
clustering_wrapper.generate_observation_fact(output_observation_fact)

df = pd.DataFrame(list(zip(dummy_to_patient.index.values.tolist(), dummy_to_patient.values.tolist())), columns=['dummy', 'patient'])
for col in df.columns:
    df[col] = df[col].astype(str)
df.to_csv(output_dummy_to_patient, quoting=csv.QUOTE_NONNUMERIC, index=False)




