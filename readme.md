# Parallel Naive Bayes
## Implementation 
The parallel implementation of the naive bayes algorithm was made using parallelism in 2 of the 3 process steps. The reading and parsing of the dataset, and the calculation of the probabilities for all the classes using the bayes theorem.

## Concepts used
- Bayes theorem
- Mutual exclusion
- Channels 
- Barriers (waitGroup)
- Async/Await pattern

## Steps
1. First of all we need to read and parse the dataset, so using parallelism, we divide the whole dataset in N chunks and process them in differents threads. To achieve it we need to use mutual exclusion in order to keep the data's integrity, and a barrier to wait for the processes to finish. 

2. Now that we the dataset was parsed, we need to calculate the probabilities for the whole vocabulary in relation to the classes. Each class have it's own probablities, so we divide the process in N sub-process where N is the number of classes. In each subprocess we use the bayes theorem to calculate the probabilities. This was implemented using the async/await pattern with channels to wait for the results of each subprocess.

3. With the probabilities calculated in the last step, we can multiply the values for the words in a text in relation to each class and classify it as spam o ham, depending of the mayor value.

