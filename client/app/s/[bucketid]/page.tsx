'use client';

import { useParams, useRouter } from 'next/navigation';
import { useState, useEffect } from 'react';
import Nav from '@/components/Nav';
import FileDropzone from '@/components/FileDropzone';
import FileContainer from '@/components/FileContainer';
import BucketExpiry from '@/components/BucketExpiry';
import BucketPasswordModal from '@/components/BucketPasswordModal';
import { useAuth } from '@/components/hooks/useAuth';
import { checkBucketProtected, fetchBucketFiles, bucketTokenStorage } from '@/lib/api';

export default function BucketPage() {
  const params = useParams();
  const router = useRouter();
  const bucketId = params.bucketid as string;
  const [refreshKey, setRefreshKey] = useState(0);
  const [showPasswordModal, setShowPasswordModal] = useState(false);
  const [isChecking, setIsChecking] = useState(true);
  const [bucketNotFound, setBucketNotFound] = useState(false);
  const { isAdmin } = useAuth({ bucketId });

  useEffect(() => {
    const checkProtection = async () => {
      if (!bucketId) return;

      setIsChecking(true);
      try {
        const protectionStatus = await checkBucketProtected(bucketId);
        
        if (protectionStatus.protected) {
          // Check if we have a token
          const token = bucketTokenStorage.getBucketAccessToken(bucketId);
          
          if (!token) {
            // No token, show modal
            setShowPasswordModal(true);
            setIsChecking(false);
            return;
          }

          // Verify token works by trying to fetch bucket
          try {
            await fetchBucketFiles(bucketId);
            // Token is valid
            setShowPasswordModal(false);
          } catch (error) {
            // Token is invalid or expired, clear it and show modal
            bucketTokenStorage.clearBucketAccessToken(bucketId);
            setShowPasswordModal(true);
          }
        } else {
          // Not protected, no modal needed
          setShowPasswordModal(false);
        }
      } catch (error) {
        const err = error as Error & { status?: number };
        if (err.status === 404) {
          setBucketNotFound(true);
          setIsChecking(false);
          return;
        }
        console.error('Error checking bucket protection:', error);
        setShowPasswordModal(false);
      } finally {
        setIsChecking(false);
      }
    };

    checkProtection();
  }, [bucketId]);

  const handleUploadSuccess = (id: string) => {
    // Upload API creates a new bucket, so navigate to it
    router.push(`/s/${id}`);
  };

  const handleAuthSuccess = () => {
    setShowPasswordModal(false);
    setRefreshKey(prev => prev + 1); // Refresh file container
  };

  if (isChecking) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center bg-black">
        <Nav />
        <main className="flex flex-col items-center justify-center w-full max-w-4xl px-4 py-8 pt-24">
          <div className="text-center">
            <p className="text-gray-400">Checking bucket access...</p>
          </div>
        </main>
      </div>
    );
  }

  if (bucketNotFound) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center bg-black">
        <Nav />
        <main className="flex flex-col items-center justify-center w-full max-w-4xl px-4 py-8 pt-24">
          <div className="text-center">
            <p className="text-xl font-semibold mb-4 text-gray-400">
              That bucket doesn&apos;t seem to exist. Are you sure you have the right code?
            </p>
            <button
              type="button"
              onClick={() => router.push('/')}
              className="py-2 px-6 rounded-lg font-medium transition-colors"
              style={{ backgroundColor: '#6A4A98', color: '#ffffff' }}
            >
              Go home
            </button>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-black">
      <Nav />
      {showPasswordModal && (
        <BucketPasswordModal
          bucketId={bucketId}
          onSuccess={handleAuthSuccess}
        />
      )}
      <main className="flex flex-col items-center justify-center w-full max-w-4xl px-4 py-8 pt-24">
        <div className="text-center mb-8">
          <h1 className="text-7xl font-black mb-4 tracking-tight" style={{ color: '#6A4A98' }}>
            CTHULHU
          </h1>

          <BucketExpiry bucketId={bucketId} />

          <p className="text-2xl font-semibold mb-2" style={{ color: '#6A4A98' }}>
            Bucket: {bucketId}
          </p>
          <p className="text-base text-gray-400">
            {isAdmin ? 'Upload more files or manage existing ones' : 'View files in this bucket'}
          </p>
        </div>

        {isAdmin && <FileDropzone onUploadSuccess={handleUploadSuccess} />}
        <div className='w-full' key={refreshKey}>
          <FileContainer bucketId={bucketId} />
        </div>
      </main>
    </div>
  );
}

