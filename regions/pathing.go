package regions

type PathFilter func(region *Region, eval interface{}) bool

type pathStep struct {
	path []RegionId
	pos  RegionId
}

/*
Path will return the shortest path between src and dst in self, discounting all paths that don't match the filter. A nil filter matches all nodes.
*/
func (self Regions) Path(src, dst RegionId, filter PathFilter, eval interface{}) (result []RegionId) {
	// queue of paths to try
	queue := []pathStep{
		pathStep{
			path: nil,
			pos:  src,
		},
	}
	// found shortest paths to the regions
	paths := map[RegionId][]RegionId{
		src: nil,
	}
	// next step preallocated
	step := pathStep{}
	// best path to the dest so far
	var best []RegionId
	// as long as we have new paths to try
	for len(queue) > 0 {
		// pick first path to try
		step = queue[0]
		// pop the queue
		queue = queue[1:]
		// if the region actually exists
		if region, found := self[step.pos]; found {
			// for each edge from the region
			for _, edge := range region.Edges {
				// if we either haven't been where this edge leads before, or we would get there along a shorter path this time (*1)
				if lastPathHere, found := paths[edge.Dst]; !found || len(step.path)+1 < len(lastPathHere) {
					// if we either haven't found dst yet, or if following this path is shorter than where we found dst
					if best == nil || len(step.path)+1 < len(best) {
						// if we aren't filtering region, or this region matches the filter
						if filter == nil || filter(region, eval) {
							// make a new path that is the path here + this region + the edge we want to follow
							thisPath := make([]RegionId, len(step.path)+1)
							// copy the path to here to the new path
							copy(thisPath, step.path)
							// add this region
							thisPath[len(step.path)] = edge.Dst
							// remember that this is the best way so far (guaranteed by *1)
							paths[edge.Dst] = thisPath
							// if this path leads to dst
							if edge.Dst == dst {
								best = thisPath
							}
							// queue up following this path further
							queue = append(queue, pathStep{
								path: thisPath,
								pos:  edge.Dst,
							})
						}
					}
				}
			}
		}
	}
	return paths[dst]
}
