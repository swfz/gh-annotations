# CLAUDE.md

## 概要

`gh annotations` は GitHub CLI の拡張。対象リポジトリの各 workflow の直近の run について、job ごとの annotation(warning / failure など)を集めて表形式、または `-json` で JSON として出力する。

```bash
gh annotations                      # カレントディレクトリのリポジトリ
gh annotations -repo owner/repo     # リポジトリを指定
gh annotations -json                # JSON 出力
```

## コマンド

```bash
go build                    # ビルド(バイナリ gh-annotations を生成。.gitignore 済み)
go test ./...               # テスト
go vet ./...                # 静的検査
gofmt -l .                  # 整形漏れの確認(出力が空なら OK)
golangci-lint run ./...     # lint(設定は .golangci.yml。CI の lint ジョブと同じ)
```

- テストは GitHub API を httpmock で差し替え、どのケースも `Options.repo` を明示するため、`GH_REPO` の指定も `gh` の認証も不要(未設定で通ることを確認済み)
- README のテスト手順は `GH_REPO=swfz/gh-annotations go test -v` だが、現状のテストでは `GH_REPO` は無くても通る

## 構成

| パス | 役割 |
|---|---|
| `main.go` | 全ロジック(`package main` の単一ファイル) |
| `main_test.go` | `run` のテーブルテスト。`httpmock` で REST API の応答を `fixtures/` から返す |
| `fixtures/*.json` | httpmock 用の API レスポンス(workflow runs / jobs / check run annotations) |
| `.github/workflows/ci.yml` | push(main)と PR で build / test / lint |
| `.github/workflows/release.yml` | `v*` タグで `cli/gh-extension-precompile` によりビルドしてリリース |

`main.go` の主な関数:

- `fetchWorkflowRuns` / `fetchJobs` / `fetchAnnotations`: REST API(`actions/runs` → `runs/{id}/jobs` → `check-runs/{id}/annotations`)を順に呼ぶ
- `latestRuns`: workflow ごとに最新の run だけを残す
- `toRecord`: run / job / annotation を出力用の 1 行(`Record`)にまとめる
- `run`: 上記を組み合わせて表(`tableprinter`)または JSON で出力する
- `main`: フラグ(`-repo` / `-json`)を解析して `run` を呼ぶ

## セットアップ前提

- Go のバージョンは `go.mod` に従う(CI も `go-version-file: go.mod`)
- 実際に `gh annotations` を実行するには `gh auth login` 済みであること
- 環境変数
  - `GH_TOKEN`: `gh` の認証トークン。CI のテストジョブでは `secrets.GITHUB_TOKEN` を渡している(テスト自体は httpmock なので必須ではない)
  - `GH_REPO`: `-repo` 未指定時に対象リポジトリを上書きする(go-gh の `repository.Current()` が参照する)

## リリース

`v*` タグを push すると `release.yml` が走り、`cli/gh-extension-precompile` がビルドしたバイナリをリリースに添付する。
