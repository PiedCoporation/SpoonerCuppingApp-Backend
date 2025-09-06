import axios, { AxiosError, type AxiosInstance } from "axios";
// import { AuthResponse } from "src/types/auth.type";
import {
  clearAuthDataFromAS,
  getAccessTokenToAS,
  saveAccessTokenToAS,
  setProfileToAS,
} from "@/utils/auth";
import path from "@/constants/path";
import config from "@/constants/config";
import { HttpStatusCode } from "@/constants/httpStatusCode.enum";

class Http {
  instance: AxiosInstance;
  private accessToken: string = "";
  constructor() {
    this.instance = axios.create({
      baseURL: `${config.baseURL}`,
      timeout: 10000,
      headers: {
        "Content-Type": "application/json",
      },
    });

    // Initialize access token asynchronously
    this.initializeAccessToken();

    this.instance.interceptors.request.use(
      (config) => {
        if (this.accessToken && config.headers) {
          config.headers.Authorization = this.accessToken;
          return config;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );
    // Add a response interceptor
    this.instance.interceptors.response.use(
      (response) => {
        const { url } = response.config;
        if (url === path.login || url === path.register) {
          const data = response.data as any;
          this.accessToken = (response.data as any).data.access_token;
          saveAccessTokenToAS(this.accessToken);
          setProfileToAS(data.data.user);
        } else if (url === path.logout) {
          this.accessToken = "";
          clearAuthDataFromAS();
        }
        return response;
      },
      (error: AxiosError) => {
        if (error.response?.status !== HttpStatusCode.UnprocessableEntity) {
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
          const data: any | undefined = error.response?.data;
          const message = data?.message || error.message;
        }
        if (error.response?.status === HttpStatusCode.Unauthorized) {
          clearAuthDataFromAS();
          // window.location.reload()
        }
        return Promise.reject(error);
      }
    );
  }

  /**
   * Initialize access token from AsyncStorage
   */
  private async initializeAccessToken(): Promise<void> {
    try {
      const token = await getAccessTokenToAS();
      if (token) {
        this.accessToken = token;
      }
    } catch (error) {
      console.error("Error initializing access token:", error);
    }
  }
}

const http = new Http().instance;

export default http;
