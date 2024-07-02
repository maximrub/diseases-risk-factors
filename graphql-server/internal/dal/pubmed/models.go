package pubmed

import "encoding/xml"

type ArticleSummery struct {
	XMLName       xml.Name `xml:"PubmedArticleSet"`
	Text          string   `xml:",chardata"`
	PubmedArticle struct {
		Text            string `xml:",chardata"`
		MedlineCitation struct {
			Text           string `xml:",chardata"`
			Status         string `xml:"Status,attr"`
			Owner          string `xml:"Owner,attr"`
			IndexingMethod string `xml:"IndexingMethod,attr"`
			PMID           struct {
				Text    string `xml:",chardata"`
				Version string `xml:"Version,attr"`
			} `xml:"PMID"`
			DateCompleted struct {
				Text  string `xml:",chardata"`
				Year  string `xml:"Year"`
				Month string `xml:"Month"`
				Day   string `xml:"Day"`
			} `xml:"DateCompleted"`
			DateRevised struct {
				Text  string `xml:",chardata"`
				Year  string `xml:"Year"`
				Month string `xml:"Month"`
				Day   string `xml:"Day"`
			} `xml:"DateRevised"`
			Article struct {
				Text     string `xml:",chardata"`
				PubModel string `xml:"PubModel,attr"`
				Journal  struct {
					Text string `xml:",chardata"`
					ISSN struct {
						Text     string `xml:",chardata"`
						IssnType string `xml:"IssnType,attr"`
					} `xml:"ISSN"`
					JournalIssue struct {
						Text        string `xml:",chardata"`
						CitedMedium string `xml:"CitedMedium,attr"`
						Volume      string `xml:"Volume"`
						PubDate     struct {
							Text string `xml:",chardata"`
							Year string `xml:"Year"`
						} `xml:"PubDate"`
					} `xml:"JournalIssue"`
					Title           string `xml:"Title"`
					ISOAbbreviation string `xml:"ISOAbbreviation"`
				} `xml:"Journal"`
				ArticleTitle string `xml:"ArticleTitle"`
				Pagination   struct {
					Text       string `xml:",chardata"`
					StartPage  string `xml:"StartPage"`
					MedlinePgn string `xml:"MedlinePgn"`
				} `xml:"Pagination"`
				ELocationID []struct {
					Text    string `xml:",chardata"`
					EIdType string `xml:"EIdType,attr"`
					ValidYN string `xml:"ValidYN,attr"`
				} `xml:"ELocationID"`
				Abstract struct {
					Text         string `xml:",chardata"`
					AbstractText []struct {
						Text        string `xml:",chardata"`
						Label       string `xml:"Label,attr"`
						NlmCategory string `xml:"NlmCategory,attr"`
					} `xml:"AbstractText"`
					CopyrightInformation string `xml:"CopyrightInformation"`
				} `xml:"Abstract"`
				AuthorList struct {
					Text       string `xml:",chardata"`
					CompleteYN string `xml:"CompleteYN,attr"`
					Author     []struct {
						Text            string `xml:",chardata"`
						ValidYN         string `xml:"ValidYN,attr"`
						LastName        string `xml:"LastName"`
						ForeName        string `xml:"ForeName"`
						Initials        string `xml:"Initials"`
						AffiliationInfo struct {
							Text        string `xml:",chardata"`
							Affiliation string `xml:"Affiliation"`
						} `xml:"AffiliationInfo"`
					} `xml:"Author"`
				} `xml:"AuthorList"`
				Language            string `xml:"Language"`
				PublicationTypeList struct {
					Text            string `xml:",chardata"`
					PublicationType []struct {
						Text string `xml:",chardata"`
						UI   string `xml:"UI,attr"`
					} `xml:"PublicationType"`
				} `xml:"PublicationTypeList"`
				ArticleDate struct {
					Text     string `xml:",chardata"`
					DateType string `xml:"DateType,attr"`
					Year     string `xml:"Year"`
					Month    string `xml:"Month"`
					Day      string `xml:"Day"`
				} `xml:"ArticleDate"`
			} `xml:"Article"`
			MedlineJournalInfo struct {
				Text        string `xml:",chardata"`
				Country     string `xml:"Country"`
				MedlineTA   string `xml:"MedlineTA"`
				NlmUniqueID string `xml:"NlmUniqueID"`
				ISSNLinking string `xml:"ISSNLinking"`
			} `xml:"MedlineJournalInfo"`
			CitationSubset  string `xml:"CitationSubset"`
			MeshHeadingList struct {
				Text        string `xml:",chardata"`
				MeshHeading []struct {
					Text           string `xml:",chardata"`
					DescriptorName struct {
						Text         string `xml:",chardata"`
						UI           string `xml:"UI,attr"`
						MajorTopicYN string `xml:"MajorTopicYN,attr"`
					} `xml:"DescriptorName"`
					QualifierName []struct {
						Text         string `xml:",chardata"`
						UI           string `xml:"UI,attr"`
						MajorTopicYN string `xml:"MajorTopicYN,attr"`
					} `xml:"QualifierName"`
				} `xml:"MeshHeading"`
			} `xml:"MeshHeadingList"`
			KeywordList struct {
				Text    string `xml:",chardata"`
				Owner   string `xml:"Owner,attr"`
				Keyword []struct {
					Text         string `xml:",chardata"`
					MajorTopicYN string `xml:"MajorTopicYN,attr"`
				} `xml:"Keyword"`
			} `xml:"KeywordList"`
			CoiStatement string `xml:"CoiStatement"`
		} `xml:"MedlineCitation"`
		PubmedData struct {
			Text    string `xml:",chardata"`
			History struct {
				Text          string `xml:",chardata"`
				PubMedPubDate []struct {
					Text      string `xml:",chardata"`
					PubStatus string `xml:"PubStatus,attr"`
					Year      string `xml:"Year"`
					Month     string `xml:"Month"`
					Day       string `xml:"Day"`
					Hour      string `xml:"Hour"`
					Minute    string `xml:"Minute"`
				} `xml:"PubMedPubDate"`
			} `xml:"History"`
			PublicationStatus string `xml:"PublicationStatus"`
			ArticleIdList     struct {
				Text      string `xml:",chardata"`
				ArticleId []struct {
					Text   string `xml:",chardata"`
					IdType string `xml:"IdType,attr"`
				} `xml:"ArticleId"`
			} `xml:"ArticleIdList"`
			ReferenceList struct {
				Text      string `xml:",chardata"`
				Reference []struct {
					Text          string `xml:",chardata"`
					Citation      string `xml:"Citation"`
					ArticleIdList struct {
						Text      string `xml:",chardata"`
						ArticleId []struct {
							Text   string `xml:",chardata"`
							IdType string `xml:"IdType,attr"`
						} `xml:"ArticleId"`
					} `xml:"ArticleIdList"`
				} `xml:"Reference"`
			} `xml:"ReferenceList"`
		} `xml:"PubmedData"`
	} `xml:"PubmedArticle"`
}
