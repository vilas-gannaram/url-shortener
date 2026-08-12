import { check } from 'k6';
import http from 'k6/http';
import { SharedArray } from 'k6/data';

const BASE_URL = 'https://lynkr-tor7.onrender.com';

// Pool of real shortcodes to hit, loaded once and shared across all VUs.
// Regenerate shortcodes.jsonl from the DB before running (see README in this dir).
const shortkeys = new SharedArray('shortkeys', function () {
	return open('./shortcodes.jsonl')
		.split('\n')
		.filter(Boolean)
		.map((line) => JSON.parse(line));
});

export const options = {
	vus: 100, // 100 virtual users running simultaneously
	iterations: 10000, // Total number of requests across all VUs
};

export default function () {
	// Pick a random key from the pool
	const key = shortkeys[Math.floor(Math.random() * shortkeys.length)];

	const res = http.get(`${BASE_URL}/${key}`, { redirects: 0 });

	check(res, {
		'redirect status is 302': (r) => r.status === 302,
	});
}
