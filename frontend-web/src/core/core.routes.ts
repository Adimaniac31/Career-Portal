export type Role = "ADMIN" | "COLLEGE_ADMIN" | "STUDENT" | "GUEST"

export interface AppRoute {
	path: string
	roles: Role[]
	authRequired: boolean
}

export const ROUTES = {
	login: {
		path: "/login",
		roles: ["GUEST"],
		authRequired: false,
	},
	signup: {
		path: "/signup",
		roles: ["GUEST"],
		authRequired: false,
	},
	dashboard: {
		path: "/dashboard",
		roles: ["ADMIN", "COLLEGE_ADMIN", "STUDENT"],
		authRequired: true,
	},
	admin: {
		path: "/admin",
		roles: ["ADMIN"],
		authRequired: true,
	},
} satisfies Record<string, AppRoute>
