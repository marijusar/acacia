import { conversationService } from '@/lib/services/conversation-service';
import { NextRequest, NextResponse } from 'next/server';

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    const { conversation_id, content } = body;

    if (!conversation_id || !content) {
      return NextResponse.json(
        { message: 'Missing conversation_id or content' },
        { status: 400 }
      );
    }

    // Create a ReadableStream that the service will write to
    const stream = new ReadableStream({
      async start(controller) {
        await conversationService.sendMessageStream(
          conversation_id,
          content,
          controller
        );
      },
    });

    // Return the stream as SSE
    return new NextResponse(stream, {
      headers: {
        'Content-Type': 'text/event-stream',
        'Cache-Control': 'no-cache',
        Connection: 'keep-alive',
      },
    });
  } catch (error) {
    return NextResponse.json(
      { message: 'Internal server error' },
      { status: 500 }
    );
  }
}
