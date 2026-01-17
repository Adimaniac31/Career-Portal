export function validateCollegeEmail(email:string, domain:string) {
	if (!email || !domain) return false

	const atIndex = email.lastIndexOf("@")
	if (atIndex === -1) return false

	const emailDomain = email.slice(atIndex + 1).toLowerCase()
	return emailDomain === domain.toLowerCase()
}

