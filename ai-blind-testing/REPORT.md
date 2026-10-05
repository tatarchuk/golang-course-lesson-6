# Task 3 — AI blind test generation

> Homework — Task 3 (Lesson 6: File I/O, JSON and Testing).
> This task is **not auto-graded** by CI — it's reviewed by your mentor.
> The workflow only checks that you actually filled this file in
> (see the "Task 3 — AI report present" step in the Actions log).

## Instructions

1. Pick one function you wrote for Task 1 or Task 2 (e.g. `LoadTodos` or
   `ValidateEmail`).
2. Give an AI assistant **only the function signature** — no explanation
   of the logic, no existing test file. Ask it to generate a full
   table-driven test suite for that signature.
3. Run the AI-generated tests against your implementation.
4. Fill in the sections below.

---

## Function under test

`func ValidateEmail(s string) bool` from `validate/validate.go`.

I picked this one because it has the most "rules" and therefore the most
room for edge cases. My own table in `validate/validate_test.go` had 27
cases at the moment I ran the experiment.

## Prompt you gave the AI

I used Claude (through Claude Code) in a **fresh session that had no
access to this repository** — it could not open my implementation or my
test file, it only got the text below. I ran it after my implementation
and my own test table were finished, so that the AI could not influence
what I wrote myself.

```text
I have a Go function with this signature in package `validate`:

```go
// ValidateEmail reports whether s is a syntactically valid email address.
func ValidateEmail(s string) bool
```

Write a complete table-driven test suite for it in Go (package `validate`,
file name `validate_email_ai_test.go`), using only the standard `testing`
package (no third-party libraries). Use a slice of anonymous structs with
`name`, `input` and `want` fields and run each case with `t.Run`. Use
`t.Errorf` (not `t.Fatalf`) inside the subtests. Name the test function
`TestValidateEmail_AI`. Cover as many edge cases as you can think of (at
least 15): valid and invalid addresses, empty input, whitespace, missing
parts, dots, special characters, unicode, length limits, etc. Give me only
the Go file, no explanation.
```

The answer was a file with **55 cases**. I saved it unchanged (I only added
a comment header on top) as `validate/validate_email_ai_test.go`, so it
runs together with my own tests:

```bash
go test -v -run TestValidateEmail_AI ./validate/...
```

### Result of the first run (before I changed anything)

54 of 55 cases passed. One failed:

```text
--- FAIL: TestValidateEmail_AI/leading_hyphen_in_domain_label
    ValidateEmail("user@-example.com") = true, want false
```

## Edge cases the AI found that you had missed

1. **Leading hyphen in a domain label** (`user@-example.com`) — **a real
   bug.** My first version only checked that the domain consists of
   letters, digits, dots and hyphens; it never looked at *where* the hyphen
   is. A DNS label can't start or end with a hyphen, so I agree with the AI.
   I fixed it the TDD way: first I added `leading hyphen in domain label`
   and `trailing hyphen in domain label` to my own table, watched them
   fail, and then rewrote `validDomain` to check every label separately
   (new helper `validDomainLabel`).
2. **Domain label exactly 63 characters long** → valid. Not a bug by itself
   (it passed), but it made me realise my code had **no label length limit
   at all** — an 185-character label was accepted. The AI only tested the
   valid side of this boundary, so its suite would never have caught that.
   I added `domain label of 64 characters is too long` → `false` and the
   check `len(label) > 63`. Funny side effect: my own "exactly 254
   characters" fixture was built from one 185-character label and started
   failing, so I had to rebuild it from 63-character labels.
3. **Leading dot in the domain** (`user@.example.com`). I had a trailing
   dot in the domain and a leading dot in the local part, but not this
   combination. Passed, no bug — the per-label check now covers it
   automatically (empty first label).
4. **Underscore in the domain** (`user@my_domain.com` → invalid). Passed,
   but I had never thought about it.
5. **Tab inside the address, only spaces, trailing space.** I had a leading
   space, a space inside and a trailing newline. All passed, because I
   reject every `unicode.IsSpace` rune, but my table did not prove it.
6. **A null byte and an invalid UTF-8 sequence** (`user\x00@example.com`,
   `user@exam\xffple.com`). I had not thought about broken input at all.
   Both passed: ranging over invalid UTF-8 gives `utf8.RuneError`, which is
   not an ASCII letter, so it is rejected.
7. **Realistic "copy-pasted from somewhere" garbage**: `mailto:user@...`,
   `<user@example.com>`, `John Doe <john@example.com>`. All rejected, but
   these are exactly the inputs a real form would get, and I did not have
   a single one of them.
8. **Punycode domain** (`user@xn--80ak6aa92e.com` → valid). Interesting
   because of the double hyphen in the middle of a label — if I had been
   "clever" and banned `--` the way I ban `..`, this would have caught it.
9. More variety on the valid side: a one-character local part, digits-only
   local part, hyphen and underscore in the local part, `co.uk`, a long
   TLD, four levels of subdomains. All passed.

So: one real bug, one missing check that happened not to matter for the
AI's inputs, and a lot of holes in my table that the code handled anyway.

## Edge cases you had that the AI missed

1. **The exact "one over the limit" boundary for the total length.** I test
   254 characters → valid *and* 255 characters → invalid. The AI tested
   254 → valid and then jumped to a 324-character address for the invalid
   case. An off-by-one like `len(s) > 255` would pass the AI's suite and
   fail mine. The AI did the same thing with the domain label: 63 → valid,
   but no 64 → invalid case. I think the AI reaches for "obviously too
   long" inputs rather than the tightest boundary, because it does not
   know which limit the implementation actually uses, so it picks a value
   that is wrong under every plausible limit.
2. **Two `@` with text between them** (`student@soft@serve.academy`) — the
   AI had `first@last@example.com`, which is the same idea, so I don't
   really count this one.

Honestly, that is about it. Apart from the boundary values, the AI's list
was a superset of mine. One thing *neither* of us had: the `%` character
in the local part, which my rules explicitly allow. I added
`stu_dent%1@softserve.academy` → valid after noticing this.

## Cases where the AI's expected output was wrong

Measured against my implementation (after the hyphen fix) **none of the 55
expectations is wrong** — all pass. But several of them are not facts,
they are **design decisions the AI silently made for me**, and they only
"passed" because I happened to decide the same way:

- `user@localhost` → `false` (the AI even named the case "missing tld").
  By the RFC this is a syntactically valid address, and Go's own
  `net/mail.ParseAddress` accepts it. I reject it on purpose, but if I had
  followed `net/mail`, this AI test would be "wrong".
- `user@-example.com` → `false`. I agree, and I fixed my code because of
  it — but `net/mail.ParseAddress` accepts this address too, so it is a
  DNS-hostname rule, not an email-syntax rule. A different author could
  reasonably have called it valid.
- `user@xn--80ak6aa92e.com` → `true` but `user@exämple.com` → `false`:
  consistent with an ASCII-only validator, which is what I wrote, but
  nothing in the signature says the function is ASCII-only.

None of the cases was marked as "depends on your rules". The AI presented
every expectation with the same confidence, which is the dangerous part:
if I had implemented `ValidateEmail` differently, I would have had to go
through all 55 cases and decide for each one whether the test or the code
is wrong.

## What you'd change about your own test-writing process after this

Test both sides of every boundary. Every limit in the rules should get two
cases: the largest valid value and the smallest invalid one. In this
experiment that was missing on both sides of the table. The AI tested a
254-character address as valid and then jumped to 324 characters for the
invalid case, so an off-by-one in my length check would have passed its
whole suite. I tested only the valid side of the domain label length, which
is why I never noticed that my code had no label length check at all. The
254/255 and 63/64 pairs that are in my table now would have caught both
problems immediately, and from now on I will add such a pair for every
limit before I write the implementation.
