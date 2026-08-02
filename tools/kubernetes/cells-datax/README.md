# DataX Cells chart

`ghcr.io/sjoon99/charts/cells:0.1.5` 기준 최소 보안 패치.

- 기준 digest: `sha256:fdfc0ba1c32145edb4e6f6bfa4452b727d15e89ced402fbf506e9f78adac6a64`
- 생성 version: `0.1.6`
- 변경: Pydio 초기 관리자 비밀번호 Secret 참조

```bash
./build.sh
```

생성 파일: `dist/cells-0.1.6.tgz`

배포 전 필수 검증:

- Helm lint/template
- Secret 평문 렌더링 0건
- DataX values 렌더링
