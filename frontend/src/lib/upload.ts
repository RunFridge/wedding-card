const MAX_DIMENSION = 2048;
const JPEG_QUALITY = 0.85;
const SKIP_COMPRESSION_BYTES = 1024 * 1024;

export const UPLOAD_CONCURRENCY = 3;

export async function compressImage(file: File): Promise<File> {
  try {
    const bitmap = await createImageBitmap(file, {
      imageOrientation: 'from-image',
    });
    const scale = Math.min(
      1,
      MAX_DIMENSION / Math.max(bitmap.width, bitmap.height),
    );
    if (
      scale === 1 &&
      file.type === 'image/jpeg' &&
      file.size <= SKIP_COMPRESSION_BYTES
    ) {
      bitmap.close();
      return file;
    }

    const canvas = document.createElement('canvas');
    canvas.width = Math.round(bitmap.width * scale);
    canvas.height = Math.round(bitmap.height * scale);
    const ctx = canvas.getContext('2d');
    if (!ctx) {
      bitmap.close();
      return file;
    }
    ctx.drawImage(bitmap, 0, 0, canvas.width, canvas.height);
    bitmap.close();

    const blob = await new Promise<Blob | null>((resolve) =>
      canvas.toBlob(resolve, 'image/jpeg', JPEG_QUALITY),
    );
    if (!blob || blob.size >= file.size) return file;

    const jpegName = file.name.replace(/\.[^.]*$/, '') + '.jpg';
    return new File([blob], jpegName, { type: 'image/jpeg' });
  } catch {
    return file;
  }
}

export async function runWithConcurrency<T>(
  items: T[],
  limit: number,
  worker: (item: T) => Promise<void>,
): Promise<void> {
  let next = 0;
  const runners = Array.from(
    { length: Math.min(limit, items.length) },
    async () => {
      while (next < items.length) {
        await worker(items[next++]);
      }
    },
  );
  await Promise.all(runners);
}
