# import glob
import pandas
import scipy

reports_root = r"../caliper-workspace/reports"
outputs_root = r"../caliper-workspace/reports-csv"
test_names = ["addOperators", "requestPermits"]#, "logBeacons"]
tps_range = range(20, 601, 20)
# tps_range = [2 ** x for x in range(0, 7, 1)] + [73.5167, 90.5097, 111.43, 128, 137.187, 168.897, 207.937, 256, 315.173, 388.023, 477.713] + [2 ** x for x in range(9, 15, 1)]
# tps_range = [2 ** x for x in range(0, 15, 1)]
run_range = range(1, 4, 1)
columns = ["TargetSendRate", "Succ", "Fail", "SendRate", "MaxLatency", "MinLatency", "AvgLatency", "Throughput", "DropRatio"]
operations = ["mean", "std"]
column_names = [op + col for col in columns[1:] for op in operations]
column_names = [columns[0]] + column_names

# def mean_zscore(group, z=3):
# 	inliers = group[group.transform(scipy.stats.zscore).abs() < z]
# 	return inliers.mean()

# def std_zscore(group, z=3):
# 	inliers = group[group.transform(scipy.stats.zscore).abs() < z]
# 	return inliers.std()

# def mean_percentile(group, q=0.9):
# 	inliers = group[group.transform(lambda x : x.quantile(1-q) < x < x.quantile(q))]
# 	return inliers.mean()

# def std_percentile(group, q=0.9):
# 	inliers = group[group.transform(lambda x : x.quantile(1-q) < x < x.quantile(q))]
# 	return inliers.std()

for test_name in test_names:
	test_data = pandas.DataFrame()
	for tps in tps_range:
		for run in run_range:
			report_path = f"{reports_root}/{test_name}-{tps}-{run}.html"
			# with open(report_path) as report_file:
			try:
				row = pandas.read_html(report_path)[0]
				row["DropRatio"] = row["Fail"] / (row["Succ"] + row["Fail"])
				test_data = pandas.concat([test_data, row])
			except ValueError:
				continue
	# row["Name"][0] = row["Name"][0].split("-")[-1]
	test_data["Name"].replace({f"{test_name}-": ""}, regex=True, inplace=True)
	test_data["Name"] = test_data["Name"].astype(float)

	# groups = test_data.groupby(["Name"], axis=0, as_index=False)
	# inliers = groups.transform(scipy.stats.zscore).abs() < 3
	test_data = test_data.groupby(["Name"], as_index=False).aggregate(operations)
	# test_data.index.name = column_names[0]
	test_data.columns = column_names

	test_data.to_csv(f"{outputs_root}/{test_name}-tps-{tps_range[0]}-{tps_range[-1]}-run-{run_range[0]}-{run_range[-1]}.csv", index=False)

