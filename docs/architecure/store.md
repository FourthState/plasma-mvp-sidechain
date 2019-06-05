# Store Design #

The store package provides the backend to storage of all information necessary for the sidechain to function.
There are 2 stores, the block store and output store.

## Block Store ##
The block store maintains all necessary information related to each plasma block produced. 
The Block type within the block store wraps the tendermint block it was committed at which is necessary for querying. 
The Block store keeps a counter for the current and next plasma block number to be used. 

## Output Store ##
All deposits, fees, and regular outputs can be stored and queried from the output store. 



