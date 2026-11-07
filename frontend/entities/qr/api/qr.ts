import { baseApi } from '@/shared/lib/api';
import type { Group } from '@/entities/group';

const responseToDataUrl = async (response: Response) => {
  if (!response.ok) {
    const contentType = response.headers.get('content-type') ?? '';
    return contentType.includes('application/json') ? response.json() : response.text();
  }

  const bytes = new Uint8Array(await response.arrayBuffer());
  let binary = '';
  for (const byte of bytes) {
    binary += String.fromCharCode(byte);
  }

  const contentType = response.headers.get('content-type') ?? 'image/png';
  return `data:${contentType};base64,${btoa(binary)}`;
};

export interface ReadQrCodeResponse {
  id: string;
  group: Group;
  key: string;
}

export const qrApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getMyQrCode: builder.query<string, void>({
      query: () => ({
        url: 'qr-code/generate',
        method: 'POST',
        responseHandler: responseToDataUrl,
      }),
    }),
    readQrCode: builder.mutation<ReadQrCodeResponse, { value: string }>({
      query: ({ value }) => ({
        url: 'qr-code/read',
        method: 'POST',
        body: { value },
      }),
    }),
  }),
});

export const { useGetMyQrCodeQuery, useReadQrCodeMutation } = qrApi;
