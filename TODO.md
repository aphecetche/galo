1. bufferize file reading in ForEachEvent
1. deduplicate clusters (i.e. when getting n exactly identical, throw away n-1)
1. identify (in the remaining clusters) those which share preclusters, i.e. those that were split 
