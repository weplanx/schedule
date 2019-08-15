<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use lumen\extra\facade\Auth;

class AuthVerify
{
    protected $scene;

    /**
     * Handle an incoming request.
     *
     * @param \Illuminate\Http\Request $request
     * @param \Closure $next
     * @return mixed
     */
    public function handle(Request $request, Closure $next)
    {
        if (empty($this->scene)) {
            return $next($request);
        }

        $result = Auth::verify($this->scene);
        if ($result) {
            return response()->json([
                'error' => 1,
                'msg' => 'token invalid'
            ]);
        } else {
            $this->definedSymbol($request);
            return $next($request);
        }
    }

    protected function definedSymbol(Request $request)
    {
        $symbol = Auth::symbol($this->scene);
        $request->user = $symbol['user'];
        $request->role = $symbol['role'];
    }
}
