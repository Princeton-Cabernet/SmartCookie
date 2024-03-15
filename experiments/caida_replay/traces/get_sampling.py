import csv 
import random

with open('caida_1min_flows.csv', mode="r", newline='') as read_file:
    csv_reader = csv.reader(read_file, delimiter='\t')
    #data = list(csv.reader(read_file, delimiter='\t'))
    with open('15percent_sample.csv', 'w') as wr_file:
        sample_writer = csv.writer(wr_file, delimiter='\t', quotechar='"', quoting=csv.QUOTE_MINIMAL)
        line_count = 0
        for row in csv_reader:
            prob = random.randint(0, 99)
            if line_count == 0 or line_count == 1:
                print(row)
                sample_writer.writerow(row)
                line_count += 1
            elif prob < 15:
                #print(f'\t{row[0]} works in the {row[1]} department, and was born in {row[2]}.')
                #print(row)
                sample_writer.writerow(row)
                line_count += 1 
            else:
                line_count += 1
        print(f'Processed {line_count} lines.')
