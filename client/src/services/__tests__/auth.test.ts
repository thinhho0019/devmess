import { loginUser, registerUser, checkEmailExists, sendResetPassword, confirmResetPassword } from "../auth";

describe("auth service", () => {
  const originalFetch = global.fetch;

  afterEach(() => {
    global.fetch = originalFetch;
    jest.resetAllMocks();
  });

  test("loginUser - success", async () => {
    const mockResponse = { access_token: "token123" };
    global.fetch = jest.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve(mockResponse) } as unknown as Response));

    const res = await loginUser("a@b.com", "pass");
    expect(res.access_token).toBe("token123");
  });

  test("loginUser - failure", async () => {
    const mockResponse = { error: "Invalid" };
  global.fetch = jest.fn(() => Promise.resolve({ ok: false, json: () => Promise.resolve(mockResponse) } as unknown as Response));

    await expect(loginUser("a@b.com", "bad")).rejects.toThrow("Invalid");
  });

  test("registerUser - success", async () => {
    const mockResponse = { message: "ok" };
  global.fetch = jest.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve(mockResponse) } as unknown as Response));

    const res = await registerUser("name", "a@b.com", "pass123A");
    expect(res.message).toBe("ok");
  });

  test("registerUser - failure", async () => {
    const mockResponse = { error: "exists" };
  global.fetch = jest.fn(() => Promise.resolve({ ok: false, json: () => Promise.resolve(mockResponse) } as unknown as Response));

    await expect(registerUser("name", "a@b.com", "pass")).rejects.toThrow("exists");
  });

  test("checkEmailExists - returns true when ok", async () => {
    const mockResponse = { message: "user exists" };
  global.fetch = jest.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve(mockResponse) } as unknown as Response));

    const res = await checkEmailExists("a@b.com");
    expect(res).toBe(true);
  });

  test("checkEmailExists - returns false when user not exists message on not ok", async () => {
    const mockResponse = { message: "user not exists" };
  global.fetch = jest.fn(() => Promise.resolve({ ok: false, json: () => Promise.resolve(mockResponse) } as unknown as Response));

    const res = await checkEmailExists("a@b.com");
    expect(res).toBe(false);
  });

  test("sendResetPassword - success", async () => {
    const mockResponse = { message: "sent" };
  global.fetch = jest.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve(mockResponse) } as unknown as Response));

    const res = await sendResetPassword("a@b.com");
    expect(res.message).toBe("sent");
  });

  test("sendResetPassword - failure", async () => {
    const mockResponse = { error: "fail" };
  global.fetch = jest.fn(() => Promise.resolve({ ok: false, json: () => Promise.resolve(mockResponse) } as unknown as Response));

    await expect(sendResetPassword("a@b.com")).rejects.toThrow("fail");
  });

  test("confirmResetPassword - success", async () => {
    const mockResponse = { message: "ok" };
  global.fetch = jest.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve(mockResponse) } as unknown as Response));

    const res = await confirmResetPassword("token", "newpass");
    expect(res.message).toBe("ok");
  });

  test("confirmResetPassword - failure", async () => {
    const mockResponse = { error: "bad" };
  global.fetch = jest.fn(() => Promise.resolve({ ok: false, json: () => Promise.resolve(mockResponse) } as unknown as Response));

    await expect(confirmResetPassword("token", "newpass")).rejects.toThrow("bad");
  });
});
