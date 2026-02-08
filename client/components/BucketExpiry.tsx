'use client';

import { useState, useEffect } from 'react';
import { fetchBucketLifecycle } from '@/lib/api';

interface BucketExpiryProps {
  bucketId: string;
}

function formatRelative(expiresAt: Date): string {
  const now = new Date();
  const ms = expiresAt.getTime() - now.getTime();
  if (ms <= 0) return 'Expired (auto-delete imminent)';
  const days = Math.floor(ms / (24 * 60 * 60 * 1000));
  const hours = Math.floor((ms % (24 * 60 * 60 * 1000)) / (60 * 60 * 1000));
  if (days > 0) return `Expires in ${days} day${days !== 1 ? 's' : ''}`;
  if (hours > 0) return `Expires in ${hours} hour${hours !== 1 ? 's' : ''}`;
  const minutes = Math.floor((ms % (60 * 60 * 1000)) / (60 * 1000));
  return `Expires in ${minutes} minute${minutes !== 1 ? 's' : ''}`;
}

export default function BucketExpiry({ bucketId }: BucketExpiryProps) {
  const [expiresAt, setExpiresAt] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  useEffect(() => {
    if (!bucketId) {
      setLoading(false);
      return;
    }

    let cancelled = false;
    setLoading(true);
    setError(false);

    fetchBucketLifecycle(bucketId)
      .then((data) => {
        if (!cancelled) {
          setExpiresAt(data.expires_at);
        }
      })
      .catch(() => {
        if (!cancelled) setError(true);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });

    return () => {
      cancelled = true;
    };
  }, [bucketId]);

  if (loading) {
    return (
      <p className="text-sm text-gray-500 mb-2">
        Loading expiry…
      </p>
    );
  }

  if (error || expiresAt === null) {
    return (
      <p className="text-sm text-gray-500 mb-2">
        Expiry unknown
      </p>
    );
  }

  const date = new Date(expiresAt);
  const localTime = date.toLocaleString();
  const relative = formatRelative(date);

  return (
    <p className="text-sm text-gray-400 mb-2">
      {relative} — {localTime}
    </p>
  );
}
