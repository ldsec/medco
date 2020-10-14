# Concept Clustering and Dummy Generation
In this project we tried some clustering techniques and then we implemented an algorithm for dummy generation
Code files:
  - clustering.py: contains the class that stores the clustering and dummies information. It provides the API for running the algorithms for clustering and generating the dummies. Only methods in this file should be called directly.
  - dummygen.py: contains the support functions for dummy generation
  - utils.py: contains support functions for the clustering
  - generate_csv.py: contains support functions for generating the new patient dimension and observation fact

# Usage
clustering.py contains the main class: HierarchicalClustering.
Two main steps have to be performed: clustering of concept codes and generation of dummies. See "using_clustering.ipynb" for an example of usage.
