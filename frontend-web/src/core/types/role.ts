export const Roles = {
	admin: "admin",
	college_admin: "college_admin",
	student: "student",
} as const

export type Role = typeof Roles[keyof typeof Roles]
