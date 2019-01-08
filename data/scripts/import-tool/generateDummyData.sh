#!/usr/bin/env bash
# $1 and $2 are respectively the path to the original patient_dimension and observation_fact tables
python using_clustering.py $1 $2
# $3 is the path to the dataset (where we store the dummy_to_patient.csv)
cp dummy_to_patient.csv $3