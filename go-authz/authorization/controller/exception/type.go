package exception

type BadRequestException struct {
	message string
	code    int
}

func (e BadRequestException) Error() string {
	return e.message
}

func (e BadRequestException) Code() int {
	return e.code
}

func NewBadRequestException(message string) BadRequestException {
	return BadRequestException{message: message, code: 400}
}

type ForbiddenException struct {
	message string
	code    int
}

func (e ForbiddenException) Error() string {
	return e.message
}

func (e ForbiddenException) Code() int {
	return e.code
}

func NewForbiddenException(message string) ForbiddenException {
	return ForbiddenException{message: message, code: 403}
}

type NotFoundException struct {
	message string
	code    int
}

func (e NotFoundException) Error() string {
	return e.message
}

func (e NotFoundException) Code() int {
	return e.code
}

func NewNotFoundException(message string) NotFoundException {
	return NotFoundException{message: message, code: 404}
}

type UnauthorizedException struct {
	message string
	code    int
}

func (e UnauthorizedException) Error() string {
	return e.message
}

func (e UnauthorizedException) Code() int {
	return e.code
}

func NewUnauthorizedException(message string) UnauthorizedException {
	return UnauthorizedException{message: message, code: 401}
}

type BadGatewayException struct {
	message string
	code    int
}

func (e BadGatewayException) Error() string {
	return e.message
}

func (e BadGatewayException) Code() int {
	return e.code
}

func NewBadGatewayException(message string) BadGatewayException {
	return BadGatewayException{message: message, code: 502}
}

type ConflictException struct {
	message string
	code    int
}

func (e ConflictException) Error() string {
	return e.message
}

func (e ConflictException) Code() int {
	return e.code
}

func NewConflictException(message string) ConflictException {
	return ConflictException{message: message, code: 409}
}
