import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react"

export const api = createApi({
	reducerPath: "api",
	baseQuery: fetchBaseQuery({
		baseUrl: import.meta.env.VITE_BACKEND_URL,
		credentials: "include", // VERY IMPORTANT
	}),
	endpoints: () => ({}),
})
