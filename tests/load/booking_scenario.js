import http from 'k6/http'
import { check, sleep} from 'k6'

export const options = {
    stages: [
        {duration: '10s', target: 10},
        {duration: '30s', target: 50},
        {duration: '20s', target: 100},
        {duration: '10s', target: 0}
    ], 

    thresholds: {
        http_req_duration: ['p(95) < 200'],
        http_req_failed: ['rate < 0.01'],
    }
}


export default function () {
    const payload = JSON.stringify({
        event_id: '550e8400-e29b-41d4-a716-446655440000',
    })

    const params = {
        headers: {
            'Content-Type': 'application/json',
        }
    }

    const res = http.post('http://localhost:8080/events/550e8400-e29b-41d4-a716-446655440000/bookings', payload, params)

    check(res, {
        'status is 201': (r) => r.status === 201
    })

    sleep(1)
}