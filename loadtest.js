import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  scenarios: {
    test: {
      executor: 'constant-vus',
      vus: 10,
      duration: '60s',
      exec: 'test',
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8081/api';

function randomPhone() {
  const num = Math.floor(100000000 + Math.random() * 900000000);
  return `+905${num}`;
}

export function listSent() {
  const res = http.get(`${BASE_URL}/messages/sent?limit=100&offset=0`);
  check(res, {
    'status 200': (r) => r.status === 200,
    'is json': (r) => r.headers['Content-Type']?.includes('application/json'),
  });
  sleep(1);
}

export function startStop() {
  const start = http.post(`${BASE_URL}/start-listener`);
  check(start, { 'start ok': (r) => r.status === 200 || r.status === 409 });
  sleep(1);

  const stop = http.post(`${BASE_URL}/stop-listener`);
  check(stop, { 'stop ok': (r) => r.status === 200 || r.status === 409 });
  sleep(1);
}

export function test() {
  const payload = {
    phone_number: randomPhone(),
    content: `Test message ${__ITER}`,
  };
  const headers = { 'Content-Type': 'application/json' };

  let postRes = http.post(`${BASE_URL}/messages`, JSON.stringify(payload), { headers });
  check(postRes, {
    'create message 200': (r) => r.status === 200 || r.status === 201,
  });

  sleep(6);

  let sentRes = http.get(`${BASE_URL}/messages/sent?limit=100&offset=0`);
  check(sentRes, {
    'sent list ok': (r) => r.status === 200,
    'contains my message': (r) => r.body.includes(payload.content),
  });

  sleep(1);
}
