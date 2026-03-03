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

export function setup() {
    const uniqEmail = `${Math.random()}@loadtest.com`
    const registerPayload = JSON.stringify({
        email: uniqEmail,
        password: 'password',
    })

    const registerParams = {
        headers: {
            'Content-Type': 'application/json',
        }
    }

    const res = http.post('http://localhost:8080/auth/register', registerPayload, registerParams)

    check(res, {
        'status is 201': (r) => r.status === 201
    })

    const loginPayload = JSON.stringify({
        email: uniqEmail,
        password: 'password',
    })

    const loginParams = {
        headers: {
            'Content-Type': 'application/json',
        }
    }

    const loginRes = http.post('http://localhost:8080/auth/login', loginPayload, loginParams)

    check(loginRes, {
        'status is 200': (r) => r.status === 200
    })

    return {
        token: loginRes.json('token'),
        email: uniqEmail,
    }
}

export default function (data) {
    const payload = JSON.stringify({
        userEmail: data.email
    })
    const params = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${data.token}`,
        }
    }

    const res = http.post('http://localhost:8080/events/a0000000-0000-0000-0000-000000000001/bookings', payload, params)

    check(res, {
        'status is 201': (r) => r.status === 201
    })

    sleep(1)
}

