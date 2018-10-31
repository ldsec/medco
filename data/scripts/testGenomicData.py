# Import pandas 
import pandas as pd

# specify input data files
cli_in_fn = "../genomic/tcga_cbio/clinical_data.csv"
mut_in_fn = "../genomic/tcga_cbio/mutation_data.csv"

# reading csv file  
df_cli = pd.read_csv(cli_in_fn, sep='\t', dtype=object)
df_mut = pd.read_csv(mut_in_fn, sep='\t', dtype=object)

braf_entries = df_mut.loc[df_mut['HUGO_GENE_SYMBOL'] == 'BRAF']
braf_entries_2 = df_mut.loc[df_mut['PROTEIN_CHANGE'] == 'E600K']

a = braf_entries_1.loc[braf_entries_1['PATIENT_ID'].isin(braf_entries_2['PATIENT_ID'])]

print(len(a))