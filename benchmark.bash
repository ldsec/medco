#Launch this script in the directory where it is located if you want no surprises. Or at least from where analyzeBenchmarkResult.py is located
totalBenchmarks=100

logDir="$1" #the directory where the log files will be saved.
customLogFileName="$2" #some string that will help the log files name stand out. Have fun there.


for i in $(seq 1 10); do
	concept=/E2ETEST/e2etest/bench${i}/

	number_children=$(($i * 5))
	suffix="${customLogFileName}${number_children}"
	outputFile="${logDir}/benchmarkOutput_${suffix}"

	for i in $(seq 1 $totalBenchmarks); do
		sudo docker-compose -f docker-compose.tools.yml run medco-cli-client --user test_explore_count_global --password test concept-children "${concept}"  >> "${outputFile}"

		if [[ $(($i%5)) == 0 ]]; then
			echo "${i}/${totalBenchmarks} benchmarks done"
		fi
	done

	measurementFile="${logDir}/benchmarksMeasurements_${suffix}.csv"
	#keep only the measurement of time to execute the request that follows the grepped string.
	cat "${outputFile}" | grep "Time to execute search concept children" | grep -o "[0-9]\+\.[0-9]\+m\?s" >> "${measurementFile}"

	#TODO execute python script to analyze the result with pandas this will simply create a box plot with pandas
	#for the moment we have a python notebook to create the boxplot
	python3 analyzeBenchmarkResult.py "${measurementFile}" "${suffix}"
done