// SearchDto.ts
export class SearchDto {
  constructor(public index: string, public result: any) {}
}

// SearchResultDto.ts
export class SearchResultDto {
  constructor(public data: SearchDto[]) {}
}
