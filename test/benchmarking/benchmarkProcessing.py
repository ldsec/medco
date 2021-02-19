#python 3 script to plot explore-stats benchmarks

#we want to plot data coming from a csv with such columns type,timer_description,duration_milliseconds
import glob
import re
import pandas as pd
import numpy as np
import seaborn as sns
import matplotlib.pyplot as plt

timerFiles = glob.glob("experiments/timers*.csv")

#find all possible buckets sizes and experiment indices
regex_timer_file = re.compile('.*timers_([\d]+)_([\d]+).csv')

bucket_sizes_set = set()
experiment_indices_set = set()

global_df = pd.DataFrame()
for f in timerFiles: 
    match = regex_timer_file.match(f)
    
    #gather all lines from the csv file
    file_name = match.group(0)
    df = pd.read_csv(file_name)
    df["bucket_size"] = int(match.group(1))
    df["experiment_index"] = int(match.group(2))
    global_df = global_df.append(df)


#keep only rows which you want to plot 
client_df = global_df[global_df["type"] == "client"]
client_df["duration_milliseconds"] = pd.to_numeric(client_df["duration_milliseconds"])
print(client_df)

for description, label in [("medco-connector-decryptions", "decryption"), ("medco-connector-explore-statistics-query-remote-execution", "remote execution")]:
    description_df = client_df[client_df["timer_description"] == description]
    ax = sns.regplot(x="bucket_size", y="duration_milliseconds", data=description_df, x_estimator=np.mean)
    ax.set(ylabel = "Execution time (ms)", xlabel = "number of buckets")
    ax.set_title("Execution time for: "+label)
    plt.savefig("explore_stats_"+label+"_regression_plot.png")
    plt.show()

#building the plot