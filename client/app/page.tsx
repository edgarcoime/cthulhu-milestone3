'use client';

import { useRouter } from 'next/navigation';
import Nav from '@/components/Nav';
import FileDropzone from '@/components/FileDropzone';

export default function Home() {
  const router = useRouter();

  const handleUploadSuccess = (id: string) => {
    router.push(`/s/${id}`);
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-black">
      <Nav />
      <main className="flex flex-col items-center justify-center w-full max-w-4xl px-4 py-8 pt-24">
        <div className="text-center mb-8">
          <h1 className="text-7xl font-black mb-4 tracking-tight" style={{ color: '#6A4A98' }}>
            CTHULHU
          </h1>
          <p className="text-2xl font-semibold mb-2" style={{ color: '#6A4A98' }}>
            Sharing files anonymously
          </p>
          <p className="text-base text-gray-400">
            no accounts no tracking just upload and share
          </p>
        </div>

        <FileDropzone onUploadSuccess={handleUploadSuccess} />
      </main>
    </div>
  );
}
