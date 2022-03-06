package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestSampleInput(t *testing.T) {
	input := "Many hands make light work."
	expected := "TWFueSBoYW5kcyBtYWtlIGxpZ2h0IHdvcmsu"
	var buffer bytes.Buffer
	btoa(strings.NewReader(input), &buffer)
	actual := buffer.String()
	if actual != expected {
		t.Logf("\nExpected:\n%v\nActual:\n%v", expected, actual)
		t.Fail()
	}
}

func TestSampleBigInput(t *testing.T) {
	input := "On the other hand, we denounce with righteous indignation and dislike men who are so beguiled and demoralized by the charms of pleasure of the moment, so blinded by desire, that they cannot foresee the pain and trouble that are bound to ensue; and equal blame belongs to those who fail in their duty through weakness of will, which is the same as saying through shrinking from toil and pain. These cases are perfectly simple and easy to distinguish. In a free hour, when our power of choice is untrammelled and when nothing prevents our being able to do what we like best, every pleasure is to be welcomed and every pain avoided. But in certain circumstances and owing to the claims of duty or the obligations of business it will frequently occur that pleasures have to be repudiated and annoyances accepted. The wise man therefore always holds in these matters to this principle of selection: he rejects pleasures to secure other greater pleasures, or else he endures pains to avoid worse pains.."
	expected := "T24gdGhlIG90aGVyIGhhbmQsIHdlIGRlbm91bmNlIHdpdGggcmlnaHRlb3VzIGluZGlnbmF0aW9uIGFuZCBkaXNsaWtlIG1lbiB3aG8gYXJlIHNvIGJlZ3VpbGVkIGFuZCBkZW1vcmFsaXplZCBieSB0aGUgY2hhcm1zIG9mIHBsZWFzdXJlIG9mIHRoZSBtb21lbnQsIHNvIGJsaW5kZWQgYnkgZGVzaXJlLCB0aGF0IHRoZXkgY2Fubm90IGZvcmVzZWUgdGhlIHBhaW4gYW5kIHRyb3VibGUgdGhhdCBhcmUgYm91bmQgdG8gZW5zdWU7IGFuZCBlcXVhbCBibGFtZSBiZWxvbmdzIHRvIHRob3NlIHdobyBmYWlsIGluIHRoZWlyIGR1dHkgdGhyb3VnaCB3ZWFrbmVzcyBvZiB3aWxsLCB3aGljaCBpcyB0aGUgc2FtZSBhcyBzYXlpbmcgdGhyb3VnaCBzaHJpbmtpbmcgZnJvbSB0b2lsIGFuZCBwYWluLiBUaGVzZSBjYXNlcyBhcmUgcGVyZmVjdGx5IHNpbXBsZSBhbmQgZWFzeSB0byBkaXN0aW5ndWlzaC4gSW4gYSBmcmVlIGhvdXIsIHdoZW4gb3VyIHBvd2VyIG9mIGNob2ljZSBpcyB1bnRyYW1tZWxsZWQgYW5kIHdoZW4gbm90aGluZyBwcmV2ZW50cyBvdXIgYmVpbmcgYWJsZSB0byBkbyB3aGF0IHdlIGxpa2UgYmVzdCwgZXZlcnkgcGxlYXN1cmUgaXMgdG8gYmUgd2VsY29tZWQgYW5kIGV2ZXJ5IHBhaW4gYXZvaWRlZC4gQnV0IGluIGNlcnRhaW4gY2lyY3Vtc3RhbmNlcyBhbmQgb3dpbmcgdG8gdGhlIGNsYWltcyBvZiBkdXR5IG9yIHRoZSBvYmxpZ2F0aW9ucyBvZiBidXNpbmVzcyBpdCB3aWxsIGZyZXF1ZW50bHkgb2NjdXIgdGhhdCBwbGVhc3VyZXMgaGF2ZSB0byBiZSByZXB1ZGlhdGVkIGFuZCBhbm5veWFuY2VzIGFjY2VwdGVkLiBUaGUgd2lzZSBtYW4gdGhlcmVmb3JlIGFsd2F5cyBob2xkcyBpbiB0aGVzZSBtYXR0ZXJzIHRvIHRoaXMgcHJpbmNpcGxlIG9mIHNlbGVjdGlvbjogaGUgcmVqZWN0cyBwbGVhc3VyZXMgdG8gc2VjdXJlIG90aGVyIGdyZWF0ZXIgcGxlYXN1cmVzLCBvciBlbHNlIGhlIGVuZHVyZXMgcGFpbnMgdG8gYXZvaWQgd29yc2UgcGFpbnMuLg=="
	var buffer bytes.Buffer
	btoa(strings.NewReader(input), &buffer)
	actual := buffer.String()
	if actual != expected {
		t.Logf("\nExpected:\n%v\nActual:\n%v", expected, actual)
		t.Fail()
	}
}

func TestSampleInputParallel(t *testing.T) {
	input := "Many hands make light work."
	expected := "TWFueSBoYW5kcyBtYWtlIGxpZ2h0IHdvcmsu"
	var buffer bytes.Buffer
	btoa_parallel(strings.NewReader(input), &buffer)
	actual := buffer.String()
	if actual != expected {
		t.Logf("\nExpected:\n%v\nActual:\n%v", expected, actual)
		t.Fail()
	}
}

func TestSampleBigInputParallel(t *testing.T) {
	input := "On the other hand, we denounce with righteous indignation and dislike men who are so beguiled and demoralized by the charms of pleasure of the moment, so blinded by desire, that they cannot foresee the pain and trouble that are bound to ensue; and equal blame belongs to those who fail in their duty through weakness of will, which is the same as saying through shrinking from toil and pain. These cases are perfectly simple and easy to distinguish. In a free hour, when our power of choice is untrammelled and when nothing prevents our being able to do what we like best, every pleasure is to be welcomed and every pain avoided. But in certain circumstances and owing to the claims of duty or the obligations of business it will frequently occur that pleasures have to be repudiated and annoyances accepted. The wise man therefore always holds in these matters to this principle of selection: he rejects pleasures to secure other greater pleasures, or else he endures pains to avoid worse pains.."
	expected := "T24gdGhlIG90aGVyIGhhbmQsIHdlIGRlbm91bmNlIHdpdGggcmlnaHRlb3VzIGluZGlnbmF0aW9uIGFuZCBkaXNsaWtlIG1lbiB3aG8gYXJlIHNvIGJlZ3VpbGVkIGFuZCBkZW1vcmFsaXplZCBieSB0aGUgY2hhcm1zIG9mIHBsZWFzdXJlIG9mIHRoZSBtb21lbnQsIHNvIGJsaW5kZWQgYnkgZGVzaXJlLCB0aGF0IHRoZXkgY2Fubm90IGZvcmVzZWUgdGhlIHBhaW4gYW5kIHRyb3VibGUgdGhhdCBhcmUgYm91bmQgdG8gZW5zdWU7IGFuZCBlcXVhbCBibGFtZSBiZWxvbmdzIHRvIHRob3NlIHdobyBmYWlsIGluIHRoZWlyIGR1dHkgdGhyb3VnaCB3ZWFrbmVzcyBvZiB3aWxsLCB3aGljaCBpcyB0aGUgc2FtZSBhcyBzYXlpbmcgdGhyb3VnaCBzaHJpbmtpbmcgZnJvbSB0b2lsIGFuZCBwYWluLiBUaGVzZSBjYXNlcyBhcmUgcGVyZmVjdGx5IHNpbXBsZSBhbmQgZWFzeSB0byBkaXN0aW5ndWlzaC4gSW4gYSBmcmVlIGhvdXIsIHdoZW4gb3VyIHBvd2VyIG9mIGNob2ljZSBpcyB1bnRyYW1tZWxsZWQgYW5kIHdoZW4gbm90aGluZyBwcmV2ZW50cyBvdXIgYmVpbmcgYWJsZSB0byBkbyB3aGF0IHdlIGxpa2UgYmVzdCwgZXZlcnkgcGxlYXN1cmUgaXMgdG8gYmUgd2VsY29tZWQgYW5kIGV2ZXJ5IHBhaW4gYXZvaWRlZC4gQnV0IGluIGNlcnRhaW4gY2lyY3Vtc3RhbmNlcyBhbmQgb3dpbmcgdG8gdGhlIGNsYWltcyBvZiBkdXR5IG9yIHRoZSBvYmxpZ2F0aW9ucyBvZiBidXNpbmVzcyBpdCB3aWxsIGZyZXF1ZW50bHkgb2NjdXIgdGhhdCBwbGVhc3VyZXMgaGF2ZSB0byBiZSByZXB1ZGlhdGVkIGFuZCBhbm5veWFuY2VzIGFjY2VwdGVkLiBUaGUgd2lzZSBtYW4gdGhlcmVmb3JlIGFsd2F5cyBob2xkcyBpbiB0aGVzZSBtYXR0ZXJzIHRvIHRoaXMgcHJpbmNpcGxlIG9mIHNlbGVjdGlvbjogaGUgcmVqZWN0cyBwbGVhc3VyZXMgdG8gc2VjdXJlIG90aGVyIGdyZWF0ZXIgcGxlYXN1cmVzLCBvciBlbHNlIGhlIGVuZHVyZXMgcGFpbnMgdG8gYXZvaWQgd29yc2UgcGFpbnMuLg=="
	var buffer bytes.Buffer
	btoa_parallel(strings.NewReader(input), &buffer)
	actual := buffer.String()
	if actual != expected {
		t.Logf("\nExpected:\n%v\nActual:\n%v", expected, actual)
		t.Fail()
	}
}
