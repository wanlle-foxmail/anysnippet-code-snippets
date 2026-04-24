# Claim Queue Jobs with SKIP LOCKED in MySQL

Claim a small batch of queued jobs in MySQL with `FOR UPDATE SKIP LOCKED` inside one transaction.

This snippet is useful when multiple workers should pull work from the same table without claiming the same queued job twice.

## Highlights

- Uses `SKIP LOCKED` for claims
- Keeps claim and update atomic
- Returns the claimed job batch

## Use Cases

- Claim the next two jobs for one worker
- Let multiple workers drain the same queue table safely
- Skip rows another worker already locked in its current transaction

## Code

```sql
START TRANSACTION;

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
```

## Notes

- This snippet intentionally claims a fixed batch of two jobs for one worker id.
- `SKIP LOCKED` only works inside a transaction that takes row locks.
- The temporary table keeps the claim set stable between the `SELECT`, `UPDATE`, and final `SELECT`.
- Add a composite queue index on `(status, priority DESC, job_id ASC)` before using this pattern on a real jobs table.
- Bind worker ids and other dynamic values with query parameters in application code instead of string interpolation.
- Verify the transaction isolation level that your queue workers use; this pattern assumes a row-locking transaction around the claim.

## Verification

Run the integration tests from the snippet root:

```bash
python -m unittest discover -s tests -p "test_*.py"
```

The verified test suite covers:

- claiming the highest-priority queued jobs
- skipping rows locked by another transaction
- limiting the batch to two jobs
- breaking priority ties by `job_id`
- ignoring non-queued jobs

The tests start a temporary MySQL 8 instance with local `mysqld`, `mysqladmin`, and `mysql` binaries.

## Files

- `src/claim_queue_jobs_with_skip_locked.sql`
- `tests/test_claim_queue_jobs_with_skip_locked.py`
- `snippet.json`