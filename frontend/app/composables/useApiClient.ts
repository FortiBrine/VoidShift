import axios, { type AxiosError, type Method, type ResponseType } from 'axios'

type ApiErrorPayload = {
  code?: string
  message?: string
  errors?: Record<string, string[]>
}

type ApiRequestOptions<TBody> = {
  method?: Method
  body?: TBody
  responseType?: ResponseType
}

const apiInstance = axios.create({
  withCredentials: true,
})

export const useApiClient = () => {
  const request = async <TResponse, TBody = unknown>(
    path: string,
    options: ApiRequestOptions<TBody> = {},
  ): Promise<TResponse> => {
    const response = await apiInstance.request<TResponse>({
      url: path,
      method: options.method,
      data: options.body,
      responseType: options.responseType,
    })

    return response.data
  }

  const getErrorPayload = (error: unknown): ApiErrorPayload | null => {
    if (!axios.isAxiosError(error)) {
      return null
    }

    const maybeData = (error as AxiosError).response?.data

    if (!maybeData || typeof maybeData !== 'object') {
      return null
    }

    return maybeData as ApiErrorPayload
  }

  const getErrorMessage = (error: unknown, fallback = 'Сталася помилка'): string => {
    const payload = getErrorPayload(error)

    if (payload?.message && payload.message.length > 0) {
      return payload.message
    }

    return fallback
  }

  return {
    request,
    getErrorPayload,
    getErrorMessage,
  }
}
