import {useEffect, useState} from 'react'
import {Api, Image as ImageInfo, Metadata} from "./Api.js";
import {AppShell, Center, Flex, Footer, Header, Image, Pagination, Text} from "@mantine/core";


function App() {
  const [images, setImages] = useState<ImageInfo[]>([])
  const [metadata, setMetadata] = useState<Metadata | null>(null)
  const [page, setPage] = useState(1)

  const api = new Api()
  const getURL = (id: string) => `${api.baseUrl}/${encodeURIComponent(id)}`

  useEffect(() => {
    const fetchImages = async () => {
      const res = await api.getImages({page: page-1})
      if (res.data.items) {
        setImages(res.data.items)
      }
      setMetadata(res.data.metadata || null)
    }
    fetchImages().catch(console.error)
  }, [page])

  const header = (
    <Header height={{base: 50, md: 70}} p="md" fixed>
      <div style={{display: 'flex', alignItems: 'center', height: '100%'}}>
        <Text>Application header</Text>
      </div>
    </Header>
  )
  const footer = (
    <Footer height={{base: 30, md: 50}} fixed>
      <Center>
        <Pagination total={metadata?.totalPages || 0} onChange={setPage}/>
      </Center>
    </Footer>
  )
  const imageList = (
    images.map((image) => (
      <Image key={image.id} src={getURL(image.id)} alt={image.prompt} radius="md" width={512} height={768}
             fit="contain" withPlaceholder/>
    ))
  )

  return (
    <AppShell padding="md" header={header} footer={footer}>
      <Flex wrap="wrap" justify="center" gap={{base: 'sm', sm: 'lg'}}>
        {imageList}
      </Flex>
    </AppShell>
  )
}

export default App
