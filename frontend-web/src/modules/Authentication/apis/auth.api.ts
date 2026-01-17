import { api } from "../../../core/api"
import type { College } from "../../../core/types/college"
import type { ApiListResponse } from "../../../core/types/api"
import type { SignupRequest } from "../Interfaces/auth.types"

export const authApi = api.injectEndpoints({
	endpoints: (builder) => ({
		getColleges: builder.query<ApiListResponse<College>, void>({
			query: () => "/api/colleges",
		}),

		signup: builder.mutation<void, SignupRequest>({
			query: (body) => ({
				url: "/api/auth/signup",
				method: "POST",
				body,
			}),
		}),
		login: builder.mutation<void, { email: string; password: string }>({
			query: (body) => ({
				url: "/api/auth/login",
				method: "POST",
				body,
				credentials: "include",
			}),
		}),

	}),
})

export const {
	useGetCollegesQuery,
	useSignupMutation,
	useLoginMutation
} = authApi

