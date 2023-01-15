import * as components from "./mporalComponents"

/**
 * @description 
 * @param req
 */
export function mpLogin(req: components.MpLoginRequest) {
	return webapi.post<components.MpLoginResponse>("/mporal/user/login", req)
}

/**
 * @description 
 */
export function wxPublicCallback() {
	return webapi.post<null>("/mporal/user/public/callback")
}

/**
 * @description 
 */
export function wxCallbackVerification() {
	return webapi.get<null>("/mporal/user/public/callback")
}

/**
 * @description 
 * @param req
 */
export function mpBindMobile(req: components.MpBindMobileRequest) {
	return webapi.post<null>("/mporal/user/bind", req)
}

/**
 * @description 
 * @param req
 */
export function mpSwitchRole(req: components.MpSwitchRoleRequest) {
	return webapi.post<null>("/mporal/user/switch/role", req)
}

/**
 * @description 
 */
export function wxPublic() {
	return webapi.get<components.WxPublicResp>("/mporal/user/public")
}

/**
 * @description 
 * @param req
 */
export function publishHomework(req: components.PublishHomeworkRequest) {
	return webapi.post<components.PublishHomeworkResponse>("/mporal/homework/publish", req)
}

/**
 * @description 
 * @param req
 */
export function publishPeriodicHomework(req: components.PublishPeriodicHomeworkRequest) {
	return webapi.post<components.PublishPeriodicHomeworkResponse>("/mporal/periodic-homework/publish", req)
}

/**
 * @description 
 * @param req
 */
export function oralReport(req: components.OralReportRequest) {
	return webapi.post<components.OralReportResponse>("/mporal/oral/report", req)
}

/**
 * @description 
 * @param req
 */
export function homeworkOralReport(req: components.HomeworkOralReportRequest) {
	return webapi.post<components.HomeworkOralReportResponse>("/mporal/homework/oral/report", req)
}

/**
 * @description 
 * @param req
 */
export function homeworkDetail(req: components.HomeworkDetailRequest) {
	return webapi.post<components.HomeworkDetailResponse>("/mporal/homework/detail", req)
}

/**
 * @description 
 * @param req
 */
export function periodicHomeworkDetail(req: components.PeriodicHomeworkDetailRequest) {
	return webapi.post<components.PeriodicHomeworkDetailResponse>("/mporal/periodic-homework/detail", req)
}

/**
 * @description 
 * @param req
 */
export function homeworkJoin(req: components.HomeworkJoinRequest) {
	return webapi.post<components.HomeworkJoinResponse>("/mporal/homework/join", req)
}

/**
 * @description 
 * @param req
 */
export function homeworkOralList(req: components.HomeworkOralListRequest) {
	return webapi.post<components.HomeworkOralListResponse>("/mporal/homework/oral/list", req)
}

/**
 * @description 
 * @param req
 */
export function homeworkOralReportUserList(req: components.HomeworkOralReportUserListRequest) {
	return webapi.post<components.HomeworkOralReportUserListResponse>("/mporal/homework/oral/report/user/list", req)
}

/**
 * @description 
 * @param req
 */
export function homeworkDataChart(req: components.HomeworkDataChartRequest) {
	return webapi.post<components.HomeworkDataChartResponse>("/mporal/homework/data-chart", req)
}

/**
 * @description 
 * @param req
 */
export function oralSave(req: components.OralSaveRequest) {
	return webapi.post<components.OralSaveResponse>("/mporal/oral/save", req)
}

/**
 * @description 
 * @param req
 */
export function oralList(req: components.OralListRequest) {
	return webapi.post<components.OralListResponse>("/mporal/oral/list", req)
}

/**
 * @description 
 * @param req
 */
export function oralDelete(req: components.OralDeleteRequest) {
	return webapi.post<components.OralDeleteResponse>("/mporal/oral/delete", req)
}

/**
 * @description 
 * @param req
 */
export function resultParse(req: components.ResultParseRequest) {
	return webapi.post<components.ResultParseResponse>("/mporal/oral/result/parse", req)
}

/**
 * @description 
 * @param req
 */
export function userSignature(req: components.UserSignatureRequest) {
	return webapi.post<components.UserSignatureResponse>("/mporal/user/signature", req)
}

/**
 * @description 
 * @param params
 */
export function submitCount(params: components.SubmitCountRequestParams) {
	return webapi.get<components.SubmitCountResponse>("/mporal/homework/submit-count/:homeworkId", params)
}

/**
 * @description 
 * @param req
 */
export function undo(req: components.UndoRequest) {
	return webapi.post<components.UndoResponse>("/mporal/homework/undo", req)
}

/**
 * @description 
 * @param req
 */
export function experience(req: components.ExperienceRequest) {
	return webapi.post<components.ExperienceResponse>("/mporal/experience/parse", req)
}

/**
 * @description 
 */
export function experienceArticle() {
	return webapi.get<components.ExperienceArticleResponse>("/mporal/experience/article")
}

/**
 * @description 
 * @param req
 */
export function roleList(req: components.RoleListRequest) {
	return webapi.post<components.RoleListResponse>("/mporal/role/homework/list", req)
}

/**
 * @description 
 */
export function showStrip() {
	return webapi.get<components.ShowStripResponse>("/mporal/showstrip")
}

/**
 * @description 
 * @param req
 */
export function getwxacode(req: components.GetwxacodeRequest) {
	return webapi.post<components.GetwxacodeResponse>("/mporal/getwxacode", req)
}

/**
 * @description 
 * @param params
 */
export function getPoster(params: components.GetPosterRequestParams) {
	return webapi.get<components.GetPosterResponse>("/mporal/poster/param/:userId", params)
}

/**
 * @description 
 * @param params
 */
export function getBookUrl(params: components.GetBookUrlRequestParams) {
	return webapi.get<components.GetBookUrlResponse>("/mporal/book/url/:bookId", params)
}
