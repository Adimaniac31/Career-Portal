import type { ApiErrorResponse } from "../../../core/types/api"
export interface SignupRequest {
	name: string
	email: string
	password: string
	college_id: number
}

export type AuthError = ApiErrorResponse
