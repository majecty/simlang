# LLVM IR Generator Sample

이 프로젝트는 간단한 Lisp 스타일 표현식을 LLVM IR로 컴파일하는 샘플 프로그램입니다.

## 기능

- Let 표현식 지원: `(let (x 10) in x)`
- 산술 연산 지원: `(+ 1 2 3)`
- 변수 바인딩 및 참조
- LLVM IR 코드 생성

## 사용법

```bash
go run .
```

프로그램은 하드코딩된 표현식을 파싱하고 LLVM IR을 생성하여 `output.ll` 파일에 저장합니다.

## 생성된 파일

- `output.ll`: 생성된 LLVM IR 코드

## 예시

입력: `(let (x 10) in x)`
출력: x 변수에 10을 바인딩하고 그 값을 반환하는 LLVM IR 코드
