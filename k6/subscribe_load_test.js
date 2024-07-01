import http from 'k6/http';
import {check, sleep} from 'k6';

export const options = {
    stages: [
        {
            duration: '10s',
            target: 100
        },
        {
            duration: '30s',
            target: 100
        },
        {
            duration: '10s',
            target: 0
        }
    ],
    thresholds: {
        http_req_duration: ['p(95)<200']
    }
};

export default function () {
    const email = `testuser${__VU}_${Math.floor(Math.random() * 10)}@example.com`;
    const res = http.post(`http://localhost:8080/api/subscribe?email=${email}`);
    check(res, {"status was 200": (r) => r.status == 200});
    sleep(1);
}