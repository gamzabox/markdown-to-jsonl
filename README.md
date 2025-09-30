# Markdown to JSON Converter (mtoj)

## 개요
mtojl은 마크다운(markdown) 문서를 RAG 임베딩에 최적화된 JSON 형식으로 변환하는 커맨드라인 프로그램입니다. Go 1.25.1로 개발되었으며, 마크다운의 헤딩, 리스트, 코드블록 등의 계층 구조를 보존하여 JSON로 변환합니다.

## 설치 및 실행
1. Go 1.25.1 이상이 설치되어 있어야 합니다.
2. 소스 코드를 클론하거나 다운로드합니다.
3. 터미널에서 프로젝트 디렉토리로 이동합니다.
4. 다음 명령어로 테스트 케이스를 실행 합니다.
   ```
   go test -v ./...
   ```
5. 다음 명령어로 빌드합니다:
   ```
   go build -o mtojl .
   ```
6. 변환할 마크다운 파일이 있는 위치에서 다음과 같이 실행합니다:
   ```
   ./mtojl example.md
   ```
   - Windows 환경에서는 `mtojl.exe`로 실행합니다.


## 사용법
```
Markdown to JSON 1.0.0
Usage: mtoj markdown-file
```
- `-h` 옵션 또는 인자를 생략하면 위 사용법 메시지가 출력됩니다.
- 입력한 마크다운 파일이 존재하지 않으면 에러 메시지를 출력하고 종료합니다.
- 출력 파일은 입력 파일명에 `.json` 확장자를 붙여 생성합니다.
- 동일한 이름의 출력 파일이 이미 존재하면 `-1`, `-2` 등의 숫자 postfix를 붙여 저장합니다.

## 기능
- 마크다운 헤딩, 리스트, 코드블록 등의 계층(depth, path) 보존
- 코드블록 최적화 (현재는 기본 구현)
- 마크다운 내 JSON 이중 색인 지원 (기본 틀 구현)
- 예외 및 경계값 처리
- 테스트 주도 개발(TDD) 방식으로 구현 및 테스트 코드 포함

## 변환 예시

### 입력 마크다운 예시
```markdown
# 제목 1
일반 텍스트 내용입니다.

## 제목 2
- 리스트 아이템 1
- 리스트 아이템 2

```go
func main() {
    fmt.Println("Hello, World!")
}
```

### 출력 JSON 예시
```json
[
   {"Type":"heading","Content":"제목 1","Depth":1,"Path":["제목 1"]},
   {"Type":"text","Content":"일반 텍스트 내용입니다.","Depth":1,"Path":["제목 1"]},
   {"Type":"heading","Content":"제목 2","Depth":2,"Path":["제목 1","제목 2"]},
   {"Type":"list","Content":"리스트 아이템 1","Depth":1,"Path":["제목 1","제목 2"]},
   {"Type":"list","Content":"리스트 아이템 2","Depth":1,"Path":["제목 1","제목 2"]},
   {"Type":"codeblock","Content":"func main() {\n    fmt.Println(\"Hello, World!\")\n}\n","Depth":2,"Path":["제목 1","제목 2"]}
]
```

### JSON 이중 색인 예시
- 워크플로우나 특정 JSONL 데이터가 마크다운 내에 포함되어 있을 때, 이중 색인을 통해 별도로 인덱싱하여 검색 효율을 높일 수 있습니다.
- 예를 들어, 다음과 같은 JSONL 데이터가 마크다운에 포함되어 있다고 가정합니다:

```json
{"workflow":"example","steps":["step1","step2"]}
```

## 테스트
프로젝트 내 테스트 코드는 다음과 같이 실행합니다:
```
go test ./...
```

## 라이선스
Apache 2.0 License
