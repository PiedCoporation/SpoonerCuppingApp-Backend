import http from "@/utils/http";

const authApi = {
  registerAccount: (body: { email: string; password: string }) =>
    http.post<any>("/register", body),
  loginAccount: (body: { email: string; password: string }) =>
    http.post<any>("/login", body),
};

export default authApi;
