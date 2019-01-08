#!/usr/bin/env python
# coding: utf-8

# In[1]:

from clustering import HierarchicalClustering
import csv
import pandas as pd

observation_fact = "data/original/observation_fact.csv"
patient_dimension = "data/original/patient_dimension.csv"


# In[2]:


print("Started hierarchical clustering.")

# Creating an instance of the clustering wrapper.
# Need to specify the observation fact path and (optionally) the separator between concept_cd 
# and modifier_cd in the new observation_fact.csv file
clustering_wrapper = HierarchicalClustering(observation_fact, concept_modifier_separator="_")

# performing hiearchical clustering
# need to specify the similarity metric and linkage method
clustering_wrapper.perform_HAC_linkage('jaccard', 'average')

# performing k-means and adapting the minimum anonymity set requirements
# need to specify the maximum ratio, the minimum allowed anonymity set and the number of clusters
clustering_wrapper.perform_clustering(max_ratio=20, min_anonymity_set_allowed=1000, n_clusters=2)

print("Finished hierarchical clustering.")

# plotting the histogram and its respective filled version (without generating any dummy)
# clustering_wrapper.plot_histogram_with_dummies_theoretical()


# In[3]:


print("Started dummy generation.")

# performing dummy generation
patient_concepts_matrix_dummies, dummy_to_patient = clustering_wrapper.generate_dummies(True)


# In[4]:


# showing ratios: theoretical vs real
# print(clustering_wrapper.theoretical_ratio(), clustering_wrapper.real_ratio())

# plotting the histograms. Last plot shows the histogram filled with dummies
# clustering_wrapper.plot_histogram()
# clustering_wrapper.plot_histogram_with_dummies_theoretical()
# clustering_wrapper.plot_histogram_with_dummies_real()


# In[5]:


# generating the new patient dimension
# need to specify the old patent_dimension.csv path and the new one
clustering_wrapper.generate_patient_dimension(patient_dimension, 'patient_dimension.csv')


# In[6]:


# generating the new observation fact
# need to specify the old observation_fact.csv and the new one
clustering_wrapper.generate_observation_fact('observation_fact.csv')


# In[7]:


df = pd.DataFrame(zip(dummy_to_patient.index.values.tolist(), dummy_to_patient.values.tolist()), columns=['dummy', 'patient'])
for col in df.columns:
    df[col] = df[col].astype(str)
df.to_csv('dummy_to_patient.csv', quoting=csv.QUOTE_NONNUMERIC, index=False)


# In[ ]:




