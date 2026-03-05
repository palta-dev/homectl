import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuthStore } from '../stores/authStore';
import type {
  ConfigResponse,
  ServicesResponse,
  HealthResponse,
} from '../types';

async function fetchAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`/api${endpoint}`, {
    ...options,
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      ...options?.headers,
    },
  });
  if (!response.ok) {
    let errorMsg = `API error: ${response.status}`;
    try {
      const errorData = await response.json();
      if (errorData.error) errorMsg = errorData.error;
    } catch {
      // Ignore if not JSON
    }
    throw new Error(errorMsg);
  }
  return response.json();
}

export function useHealth() {
  return useQuery<HealthResponse>({
    queryKey: ['health'],
    queryFn: () => fetchAPI<HealthResponse>('/health'),
    refetchInterval: 30000,
    retry: 1,
  });
}

export function useConfig() {
  return useQuery<ConfigResponse>({
    queryKey: ['config'],
    queryFn: () => fetchAPI<ConfigResponse>('/config'),
    refetchInterval: 60000,
    retry: 2,
  });
}

export function useUpdateConfig() {
  const queryClient = useQueryClient();
  const { password } = useAuthStore();
  
  return useMutation({
    mutationFn: ({ settings, password: overridePassword }: { settings: Partial<ConfigResponse>; password?: string }) => 
      fetchAPI('/config', {
        method: 'PUT',
        headers: (overridePassword || password) ? { 'X-HOMECTL-AUTH': overridePassword || password || '' } : {},
        body: JSON.stringify(settings),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['config'] });
    },
  });
}

export function useLogin() {
  const { setPassword, setAuthenticated } = useAuthStore();
  return useMutation({
    mutationFn: (password: string) => 
      fetchAPI<{ message: string }>('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ password }),
      }),
    onSuccess: (_, password) => {
      setPassword(password);
      setAuthenticated(true);
    },
  });
}

export function useLogout() {
  const { logout } = useAuthStore();
  return useMutation({
    mutationFn: () => fetchAPI('/auth/logout', { method: 'POST' }),
    onSuccess: () => logout(),
  });
}

export function useAuthCheck() {
  const { setAuthenticated, logout } = useAuthStore();
  return useQuery({
    queryKey: ['auth-check'],
    queryFn: async () => {
      try {
        await fetchAPI('/config/auth-check');
        setAuthenticated(true);
        return true;
      } catch {
        logout();
        return false;
      }
    },
    retry: false,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

export function useServices() {
  return useQuery<ServicesResponse>({
    queryKey: ['services'],
    queryFn: () => fetchAPI<ServicesResponse>('/services'),
    refetchInterval: 30000,
    retry: 2,
  });
}
