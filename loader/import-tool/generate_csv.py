import pandas as pd
import csv


def generate_patient_dimension(mapping, patient_dimension_path, output_path):
    patient_dim_df = pd.read_csv(patient_dimension_path, low_memory=False)

    patient_dim_df = patient_dim_df.fillna("")
    patient_dim_df.patient_num = patient_dim_df.patient_num.astype('str')

    patient_dim_df = patient_dim_df.set_index('patient_num')
    patient_dim_df['dummy'] = 1
    patient_dim_df['original_patient'] = ""

    for key, patient in mapping.iteritems():
        row = patient_dim_df.loc[patient].copy()
        row['dummy'] = 0
        row['original_patient'] = patient
        patient_dim_df.loc[key] = row

    for col in patient_dim_df.columns:
        patient_dim_df[col] = patient_dim_df[col].astype(str)

    patient_dim_df.to_csv(output_path, quoting=csv.QUOTE_NONNUMERIC, index=False)


def generate_observation_fact(observation_fact_df, histogram, new_patient_concepts_matrix, old_patient_concepts_matrix, clusters_labels, output_path):
    observation_fact_df = observation_fact_df[observation_fact_df['concept_cd'].isin(histogram.index)]
    observation_fact_df = observation_fact_df.fillna("")
    for col in observation_fact_df.columns:
        observation_fact_df[col] = observation_fact_df[col].astype(str)

    grouped = observation_fact_df.groupby(['patient_num', 'concept_cd']).first().reset_index()[observation_fact_df.columns]

    total_rows = new_patient_concepts_matrix.sum(1).sum()

    new_df = pd.DataFrame(index=range(total_rows),columns=grouped.columns)
    new_df.loc[grouped.index] = grouped.loc[grouped.index]

    dummy_patient_concepts_matrix = new_patient_concepts_matrix[~ new_patient_concepts_matrix.index.isin(old_patient_concepts_matrix.index)]

    i = grouped.index[-1] + 1

    real_rows = new_patient_concepts_matrix[new_patient_concepts_matrix.index.isin(old_patient_concepts_matrix.index)].sum(1).sum()
    for name, patient in dummy_patient_concepts_matrix.iterrows():
        filtered_patient = patient[patient == 1]
        for code in filtered_patient.index:
            new_df.loc[i] = pd.Series({'patient_num': filtered_patient.name, 'concept_cd': code })
            i = i+1

    new_df['cluster_label'] = clusters_labels[new_df['concept_cd']].values
    new_df['cluster_label'] = new_df['cluster_label'].astype(str)

    new_df.to_csv(output_path, quoting=csv.QUOTE_NONNUMERIC, index=False)
