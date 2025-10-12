// src/services/api.ts

interface LoginResponse {
  access_token?: string;
  message?: string;
  error?: string;

}

export const loginUser = async (email: string, password: string): Promise<LoginResponse> => {
  const url = `${import.meta.env.VITE_API_URL}/login`;

  const response = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, password }),
  });

  const data: LoginResponse = await response.json();

  if (!response.ok) {
    throw new Error(data.error || "Email hoặc mật khẩu không chính xác.");
  }

  return data;
};

interface RegisterResponse {
  message?: string;
  error?: string;
}

export const registerUser = async (name: string, email: string, password: string): Promise<RegisterResponse> => {
  const url = `${import.meta.env.VITE_API_URL}/register`;

  const response = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ name, email, password }),
  });

  const data: RegisterResponse = await response.json();

  if (!response.ok) {
    throw new Error(data.error || "Đăng ký không thành công. Vui lòng thử lại.");
  }

  return data;
};

interface CheckEmailResponse {
  message: string;
}

export const checkEmailExists = async (email: string): Promise<boolean> => {
  const url = `${import.meta.env.VITE_API_URL}/check-email`;

  const response = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email }),
  });

  const data: CheckEmailResponse = await response.json();

  if (!response.ok) {
    if (data.message === "user not exists") {
      return false; // Email đã tồn tại
    }
    throw new Error(data.message || "Kiểm tra email không thành công.");
  }

  return true;
};
export async function sendResetPassword(email: string) {
  const res = await fetch(`${import.meta.env.VITE_API_URL}/forgot-password`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email }),
  });

  if (!res.ok) {
    const data = await res.json();
    throw new Error(data.error || "Không thể gửi mail reset");
  }
  return await res.json();
}// src/services/index.ts
export async function confirmResetPassword(token: string, password: string) {
  const res = await fetch(`${import.meta.env.VITE_API_URL}/reset-password`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ token, password }),
  });

  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    console.log(data);
    throw new Error(data.error || "Không thể đặt lại mật khẩu");
  }

  return res.json();
}
