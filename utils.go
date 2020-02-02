package pool
import "time"
//binary search
func (p *Pool)binSearch(t time.Time) int {
    end := len(p.conn)-1
    start := 0
    var mid int
    for start <= end {
        mid = (start+end)/2
        //t is before expiry
        if t.Before(p.conn[mid].expiry) {
            end = mid-1
        } else {
            start = mid+1
        }
    }
    return end
}
