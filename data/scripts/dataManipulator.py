import pandas as pd
import numpy as np


# PLEASE USE THIS TO CONFIGURE YOUR SCRIPT #
# ---------------------------------------- #
# specify input data files
cli_in_fn = "../genomic/tcga_cbio/clinical_data.csv"
mut_in_fn = "../genomic/tcga_cbio/mutation_data.csv"
# specify output folder
out_folder = "../genomic/tcga_cbio/manipulations/"
# ---------------------------------------- #

df_cli = pd.read_csv(cli_in_fn, sep='\t', dtype=object)
df_mut = pd.read_csv(mut_in_fn, sep='\t', dtype=object)


def save(df_cli_new, df_mut_new):
    while True:
        opt = raw_input("\nDo you wish to save data (y|n):")

        if opt not in ["y", "n"]:
            print "Wrong option!"
            continue

        if opt == "y":
            filename = raw_input("Specify filename:")
            df_cli_new.to_csv(out_folder + filename+"_clinical_data.csv", sep='\t', encoding='utf-8', index=False)
            df_mut_new.to_csv(out_folder + filename+"_mutation_data.csv", sep='\t', encoding='utf-8', index=False)
            return

        if opt == "n":
            return


def random():
    while True:
        num_str = raw_input("\nNumber of patients for the new dataset:")

        try:
            num = int(num_str)
        except ValueError:
            print("That's not an int!")
            continue

        total_patients = df_cli["PATIENT_ID"].unique()
        list_patients = np.random.choice(total_patients, num, replace=False)

        df_cli_new = df_cli.loc[df_cli["PATIENT_ID"].isin(list_patients)]
        df_mut_new = df_mut.loc[df_mut["PATIENT_ID"].isin(list_patients)]

        save(df_cli_new, df_mut_new)
        return


def specific():
    while True:
        path = raw_input("\nPlease introduce the the path to a file with a list of PATIENT_IDs to be selected from the "
                         "original dataset.")

        try:
            df = pd.read_csv(path, sep='\t', dtype=object, header=None)
        except IOError:
            print('File not found')
            continue

        list_patients = np.array(df[0])

        df_cli_new = df_cli.loc[df_cli["PATIENT_ID"].isin(list_patients)]
        df_mut_new = df_mut.loc[df_mut["PATIENT_ID"].isin(list_patients)]

        save(df_cli_new, df_mut_new)
        return


def menu():
    while True:
        print "\n#--- MENU ---#"
        print "(1) Random"
        print "(2) Specific"
        print "(3) Help"
        print "(4) Exit"

        opt = raw_input("Option:")

        if opt not in ["1", "2", "3", "4"]:
            print "Wrong option!"
            continue

        # randomly filters the patients (and respective mutations)
        if opt == "1":
            random()

        # specific filters the patients (and respective mutations) based on a file wit patient IDs
        if opt == "2":
            specific()

        # some help noob
        if opt == "3":
            print "HELP\n"
            print "----------------------------------------------------------------------------------------------------"
            print "This is a small application to manipulate the tcga_cbio data.\n\nBEFORE ANYTHING! Check if you " \
                  "have configured the script (check the top lines on how to do it)." \
                  "\n\nWe have two modes to manipulate " \
                  "the data: \n(Random) where we sample a random number of patients to build a new smaller dataset;" \
                  "\n(Specific) where we specify (through a file) the patients' IDs we want to include." \
                  "\n\nNow go an play..."
            print "----------------------------------------------------------------------------------------------------"

        # well is an exit what do you expect
        if opt == "4":
            return


menu()
