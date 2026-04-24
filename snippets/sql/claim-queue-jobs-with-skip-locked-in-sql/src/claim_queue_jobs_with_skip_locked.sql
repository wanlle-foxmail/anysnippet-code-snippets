START TRANSACTION;

-- The example worker id and timestamp are shown inline for readability.
-- In application code, replace them with bound parameters instead of
-- interpolating dynamic values into SQL text.

CREATE TEMPORARY TABLE claimed_job_ids (
    job_id INT PRIMARY KEY
);

INSERT INTO claimed_job_ids (job_id)
SELECT job_id
FROM jobs
WHERE status = 'queued'
ORDER BY priority DESC, job_id ASC
LIMIT 2
FOR UPDATE SKIP LOCKED;

UPDATE jobs
SET status = 'processing',
    worker_id = 'worker-1',
    claimed_at = '2026-04-01 09:00:00'
WHERE job_id IN (SELECT job_id FROM claimed_job_ids);

SELECT
    job_id,
    priority,
    worker_id
FROM jobs
WHERE job_id IN (SELECT job_id FROM claimed_job_ids)
ORDER BY priority DESC, job_id ASC;

DROP TEMPORARY TABLE claimed_job_ids;

COMMIT;