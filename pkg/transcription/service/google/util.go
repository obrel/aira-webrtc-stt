package google

import speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"

func getEncoding(e string) speechpb.RecognitionConfig_AudioEncoding {
	switch e {
	case "linear16":
		return speechpb.RecognitionConfig_LINEAR16
	case "flac":
		return speechpb.RecognitionConfig_FLAC
	case "ulaw":
		return speechpb.RecognitionConfig_MULAW
	case "amr":
		return speechpb.RecognitionConfig_AMR
	case "amrwb":
		return speechpb.RecognitionConfig_AMR_WB
	case "opus":
		return speechpb.RecognitionConfig_OGG_OPUS
	case "speex":
		return speechpb.RecognitionConfig_SPEEX_WITH_HEADER_BYTE
	default:
		return speechpb.RecognitionConfig_ENCODING_UNSPECIFIED
	}
}
